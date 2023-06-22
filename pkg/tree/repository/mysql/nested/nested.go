package nested

import (
	"context"
	"database/sql"
	"errors"
	"github.com/F-Amaral/tcc/constants"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/pkg/tree/domain/repositories"
	"github.com/F-Amaral/tcc/pkg/tree/repository/mysql/nested/contracts"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"moul.io/zapgorm2"
)

type nested struct {
	db *gorm.DB
}

func NewNested(config *viper.Viper, logger log.Logger) (repositories.NestedTree, error) {
	logWrap := zapgorm2.New(logger.Desugar())
	logWrap.SetAsDefault()
	db, err := gorm.Open(mysql.Open(config.GetString(constants.NestedDbDsnKey)), &gorm.Config{Logger: logWrap})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&contracts.Node{})
	if err != nil {
		return nil, err
	}

	//err = db.Use(tracing.NewPlugin(tracing.WithDBName("nested")))
	//if err != nil {
	//	return nil, err
	//}
	return &nested{
		db: db,
	}, nil
}

func (s *nested) Save(ctx context.Context, entityNode *entity.Node) apierrors.ApiError {
	node := contracts.MapFromEntity(entityNode)
	if _, err := s.save(ctx, s.db.WithContext(ctx), node); err != nil {
		return err
	}

	return nil
}

func (s *nested) GetById(ctx context.Context, nodeId string) (*entity.Node, apierrors.ApiError) {
	node, err := s.getById(ctx, s.db.WithContext(ctx), nodeId)
	if err != nil {
		return nil, err
	}
	return contracts.MapToEntity(node), nil
}

func (s *nested) GetTree(ctx context.Context, parentId string) (*entity.Node, apierrors.ApiError) {
	parent := &contracts.Node{ID: parentId}
	res := s.db.WithContext(ctx).First(parent)
	if res.Error != nil {
		return nil, s.handleGormError(res.Error)
	}

	rows, err := s.db.WithContext(ctx).Model(&contracts.Node{}).
		Where("tree_id = ? AND lft >= ? AND rgt <= ?", parent.TreeId, parent.Left, parent.Right).
		Order("lft").Rows()
	if err != nil {
		return nil, s.handleGormError(err)
	}
	defer rows.Close()

	nodeParentMap := make(map[string]*contracts.Node)
	nodeParentMap[parent.ID] = parent
	for rows.Next() {
		var node contracts.Node
		err := s.db.WithContext(ctx).ScanRows(rows, &node)
		if err != nil {
			return nil, apierrors.NewInternalServerApiError(err.Error())
		}
		if node.ParentId == nil {
			if node.ID != parent.ID {
				return nil, apierrors.NewInternalServerApiError("more that one root node found in tree")
			}
			continue
		}

		if parentNode, ok := nodeParentMap[*node.ParentId]; ok {
			parentNode.Children = append(parentNode.Children, &node)
			nodeParentMap[node.ID] = &node
		} else {
			return nil, apierrors.NewInternalServerApiError("parent node not found")
		}
	}
	return contracts.MapToEntity(parent), nil
}

func (s *nested) AppendToTree(ctx context.Context, parentId string, entityNode *entity.Node) (*entity.Node, apierrors.ApiError) {
	if len(entityNode.Children) > 0 {
		node, err := s.handleMergeTrees(ctx, parentId, entityNode)
		if err != nil {
			return nil, err
		}
		return contracts.MapToEntity(node), nil
	}
	node, err := s.handleEmptyChild(ctx, parentId, entityNode)
	if err != nil {
		return nil, err
	}
	return contracts.MapToEntity(node), nil
}

func (s *nested) handleEmptyChild(ctx context.Context, parentId string, entityNode *entity.Node) (*contracts.Node, apierrors.ApiError) {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	parentNode, err := s.getById(ctx, tx, parentId)
	if err != nil {
		return nil, err
	}

	updateRgtCommand := "UPDATE nodes SET rgt = rgt + 2 WHERE rgt >= @parent_right AND tree_id = @tree_id"
	updateLftCommand := "UPDATE nodes SET lft = lft + 2 WHERE lft > @parent_right AND tree_id = @tree_id"

	if err := tx.Exec(updateRgtCommand,
		sql.Named("parent_right", parentNode.Right),
		sql.Named("tree_id", parentNode.TreeId)).
		Error; err != nil {
		tx.Rollback()
		return nil, s.handleGormError(err)
	}
	if err := tx.Exec(updateLftCommand,
		sql.Named("parent_right", parentNode.Right),
		sql.Named("tree_id", parentNode.TreeId)).
		Error; err != nil {
		tx.Rollback()
		return nil, s.handleGormError(err)
	}

	childNode := s.prepareNodeForAdd(entityNode, parentNode)
	if _, err := s.save(ctx, tx, childNode); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, s.handleGormError(err)
	}

	parentNode.Children = append(parentNode.Children, childNode)
	return parentNode, nil
}

func (s *nested) handleMergeTrees(ctx context.Context, parentId string, entityNode *entity.Node) (*contracts.Node, apierrors.ApiError) {
	tx := s.db.Debug().WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	parentNode, err := s.getById(ctx, tx, parentId)
	if err != nil {
		return nil, err
	}

	childNode, err := s.getById(ctx, tx, entityNode.Id)
	if err != nil {
		return nil, err
	}

	prepareChildrenQuery := `UPDATE nodes SET lft = lft + @parent_rgt - 1, rgt = rgt + @parent_rgt - 1 WHERE tree_id = @tree_id;`
	if err := tx.Exec(prepareChildrenQuery,
		sql.Named("parent_rgt", parentNode.Right),
		sql.Named("tree_id", childNode.TreeId)).
		Error; err != nil {
		tx.Rollback()
		return nil, s.handleGormError(err)
	}

	newLftOffset := childNode.Left + parentNode.Right - 1
	newRgtOffset := childNode.Right + parentNode.Right - 1

	updateParentTreeRgt := "UPDATE nodes SET rgt = rgt + @child_rgt WHERE rgt >= @parent_rgt AND tree_id = @tree_id"
	updateParentTreeLft := "UPDATE nodes SET lft = lft + @child_rgt WHERE lft > @parent_rgt AND tree_id = @tree_id"

	if err := tx.Exec(updateParentTreeRgt,
		sql.Named("child_rgt", childNode.Right),
		sql.Named("parent_rgt", parentNode.Right),
		sql.Named("tree_id", parentNode.TreeId)).
		Error; err != nil {
		tx.Rollback()
		return nil, s.handleGormError(err)
	}

	if err := tx.Exec(updateParentTreeLft,
		sql.Named("child_rgt", childNode.Right),
		sql.Named("parent_rgt", parentNode.Right),
		sql.Named("tree_id", parentNode.TreeId)).
		Error; err != nil {
		tx.Rollback()
		return nil, s.handleGormError(err)
	}

	updateChildTreeId := "UPDATE nodes SET tree_id = @parent_tree_id WHERE tree_id = @child_tree_id"
	if err := tx.Exec(updateChildTreeId,
		sql.Named("parent_tree_id", parentNode.TreeId),
		sql.Named("child_tree_id", childNode.TreeId)).
		Error; err != nil {
		tx.Rollback()
		return nil, s.handleGormError(err)
	}

	childNode.TreeId = parentNode.TreeId
	childNode.ParentId = &parentNode.ID
	childNode.Left = newLftOffset
	childNode.Right = newRgtOffset
	if _, err := s.save(ctx, tx, childNode); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	parentNode.Children = append(parentNode.Children, childNode)
	return parentNode, nil
}

func (s *nested) handleGormError(err error) apierrors.ApiError {
	if errors.Is(gorm.ErrRecordNotFound, err) {
		return apierrors.NewNotFoundApiError(err.Error())
	}
	return apierrors.NewInternalServerApiError(err.Error())
}

func (s *nested) prepareNodeForAdd(entityNode *entity.Node, parentNode *contracts.Node) *contracts.Node {
	childNode := contracts.MapFromEntity(entityNode)
	childNode.ParentId = &parentNode.ID
	childNode.TreeId = parentNode.TreeId
	childNode.Left = parentNode.Right
	childNode.Right = parentNode.Right + 1
	childNode.Level = parentNode.Level + 1
	return childNode
}

func (s *nested) save(ctx context.Context, tx *gorm.DB, relationship *contracts.Node) (*contracts.Node, apierrors.ApiError) {
	if err := tx.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(relationship).Error; err != nil {
		return nil, s.handleGormError(err)
	}

	return relationship, nil
}

func (s *nested) getById(ctx context.Context, tx *gorm.DB, nodeId string) (*contracts.Node, apierrors.ApiError) {
	node := &contracts.Node{ID: nodeId}
	if err := s.db.WithContext(ctx).Preload("Children").First(&node).Error; err != nil {
		return nil, s.handleGormError(err)
	}
	return node, nil
}

//	func (s *nested) RemoveContainerFromTree(ctx context.Context, childRelationship *entity.Node) apierrors.ApiError {
//		rb := rollback.New()
//
//		newTreeId, err := s.createNewTreeIdToRange(ctx, childRelationship.TreeID, childRelationship.GroupStart, childRelationship.GroupEnd)
//
//		if err != nil {
//			return err
//		}
//
//		rb.Add("RollbackCreateNewTreeFromRange", func() { _ = s.updateTreeId(ctx, newTreeId, childRelationship.TreeID) })
//
//		err = s.prepareChildTreeForSplit(ctx, childRelationship.GroupStart, childRelationship.GroupEnd, newTreeId)
//		if err != nil {
//			rb.Do(ctx)
//			return err
//		}
//
//		rb.Add("RollbackPreapreChildTreeForSplit", func() {
//			_ = s.rollbackChildTreeOnSplitFail(ctx, childRelationship.GroupStart, childRelationship.GroupEnd, newTreeId)
//		})
//
//		parentTreeOffset := childRelationship.GroupEnd - childRelationship.GroupStart + 1
//		err = s.prepareParentTreeForSplit(ctx, parentTreeOffset, childRelationship.TreeID)
//		if err != nil {
//			rb.Do(ctx)
//			return err
//		}
//
//		rb.Add("RollbackPrepareParentTreeForSplit", func() {
//			_ = s.rollbackParentTreeOnSplitFail(ctx, childRelationship.GroupEnd, childRelationship.TreeID)
//		})
//
//		err = s.updateParentId(ctx, childRelationship.ContainerID, "", newTreeId)
//		if err != nil {
//			rb.Do(ctx)
//			return err
//		}
//
//		return nil
//	}

//func (s *nested) createNewTreeIdToRange(ctx context.Context, oldTreeId string, rangeStart, rangeEnd uint) (string, apierrors.ApiError) {
//	query := "UPDATE container_parent SET tree_id = :new_tree_id WHERE group_start >= :child_group_start AND group_end <= :child_group_end AND tree_id = :old_tree_id;"
//
//	newTreeId, _ := uuid.NewUUID()
//
//	sqlErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//		Query(query).
//		NamedParameters(
//			sql.Named("child_group_start", rangeStart),
//			sql.Named("child_group_end", rangeEnd),
//			sql.Named("new_tree_id", newTreeId.String()),
//			sql.Named("old_tree_id", oldTreeId)).
//		Execute()
//
//	if sqlErr != nil {
//		err := errorctx.New(ctx, sqlErr.GetStatusCode(), sqlErr.GetErrorCode(), "Could not save new tree id to range")
//		errorctx.Log(ctx, err)
//		return "", err
//	}
//
//	return newTreeId.String(), nil
//}

//	func (s *nested) prepareChildTreeForSplit(ctx context.Context, rangeStart, rangeEnd uint, treeId string) apierrors.ApiError {
//		queryGroupStart := "UPDATE container_parent SET group_start = group_start - :child_group_start + 1 WHERE group_start >= :child_group_start AND group_end <= :child_group_end AND tree_id = :tree_id;"
//		queryGroupEnd := "UPDATE container_parent SET group_end = group_end - :child_group_start + 1 WHERE group_end > :child_group_start AND group_end <= :child_group_end AND tree_id = :tree_id;"
//
//		groupStartErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//			Query(queryGroupStart).
//			NamedParameters(
//				sql.Named("child_group_start", rangeStart),
//				sql.Named("child_group_end", rangeEnd),
//				sql.Named("tree_id", treeId)).
//			Execute()
//
//		if groupStartErr != nil {
//			err := errorctx.New(ctx, groupStartErr.GetStatusCode(), groupStartErr.GetErrorCode(), "Could not prepare child tree (group_start) for split")
//			errorctx.Log(ctx, err)
//			return err
//		}
//
//		groupEndErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//			Query(queryGroupEnd).
//			NamedParameters(
//				sql.Named("child_group_start", rangeStart),
//				sql.Named("child_group_end", rangeEnd),
//				sql.Named("tree_id", treeId)).
//			Execute()
//
//		if groupEndErr != nil {
//			err := errorctx.New(ctx, groupEndErr.GetStatusCode(), groupEndErr.GetErrorCode(), "Could not prepare child tree (group_end) for split")
//			errorctx.Log(ctx, err)
//			return err
//		}
//
//		return nil
//	}
//
// func (s *nested) rollbackChildTreeOnSplitFail(ctx context.Context, rangeStart, rangeEnd uint, treeId string) apierrors.ApiError {
//
//	query := "UPDATE container_parent SET group_start = group_start + @child_group_start - 1, group_end = group_end + :child_group_start - 1 WHERE tree_id = @tree_id"
//
//	sqlErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//		Query(query).
//		NamedParameters(
//			sql.Named("child_group_start", rangeStart),
//			sql.Named("child_group_end", rangeEnd),
//			sql.Named("tree_id", treeId)).
//		Execute()
//
//	if sqlErr != nil {
//		err := errorctx.New(ctx, sqlErr.GetStatusCode(), sqlErr.GetErrorCode(), "Could not rollback child tree (group_start) on split fail")
//		errorctx.Log(ctx, err)
//		return err
//	}
//
//	return nil
//
// }
//
//	func (s *nested) prepareParentTreeForSplit(ctx context.Context, childOffset uint, treeId string) apierrors.ApiError {
//		queryGroupStart := "UPDATE container_parent SET group_start = group_start - :child_group_end WHERE group_start > :child_group_end AND tree_id = @tree_id;"
//		queryGroupEnd := "UPDATE container_parent SET group_end = group_end - :child_group_end WHERE group_end > :child_group_end AND tree_id = @tree_id;"
//
//		groupStartErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//			Query(queryGroupStart).
//			NamedParameters(
//				sql.Named("child_group_end", childOffset),
//				sql.Named("tree_id", treeId)).
//			Execute()
//
//		if groupStartErr != nil {
//			err := errorctx.New(ctx, groupStartErr.GetStatusCode(), groupStartErr.GetErrorCode(), "Could not prepare parent tree (group_start) for split")
//			errorctx.Log(ctx, err)
//			return err
//		}
//
//		groupEndErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//			Query(queryGroupEnd).
//			NamedParameters(
//				sql.Named("child_group_end", childOffset),
//				sql.Named("tree_id", treeId)).
//			Execute()
//
//		if groupEndErr != nil {
//			err := errorctx.New(ctx, groupEndErr.GetStatusCode(), groupEndErr.GetErrorCode(), "Could not prepare parent tree (group_end) for split")
//			errorctx.Log(ctx, err)
//			return err
//		}
//
//		return nil
//	}
//
//	func (s *nested) rollbackParentTreeOnSplitFail(ctx context.Context, childOffset uint, treeId string) apierrors.ApiError {
//		query := "UPDATE container_parent SET group_end = group_end + :child_group_end, group_start = group_start + :child_group_end WHERE group_start > :child_group_end AND group_end > :child_group_end AND tree_id = @tree_id;"
//
//		sqlErr := s.sqlPrepareBuilder.Create(ctx, s.writer, PrefixContainer).
//			Query(query).
//			NamedParameters(
//				sql.Named("child_group_end", childOffset),
//				sql.Named("tree_id", treeId)).
//			Execute()
//
//		if sqlErr != nil {
//			err := errorctx.New(ctx, sqlErr.GetStatusCode(), sqlErr.GetErrorCode(), "Could not rollback parent tree on split fail")
//			errorctx.Log(ctx, err)
//			return err
//		}
//
//		return nil
//	}
