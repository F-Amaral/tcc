package ppt

import (
	"context"
	"errors"
	"github.com/F-Amaral/tcc/constants"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/pkg/tree/domain/repositories"
	"github.com/F-Amaral/tcc/pkg/tree/repository/mysql/ppt/contracts"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ppt struct {
	db *gorm.DB
}

func NewPpt(config *viper.Viper) (repositories.Tree, error) {
	db, err := gorm.Open(mysql.Open(config.GetString(constants.PPtDbDsnKey)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&contracts.Node{})

	//err = db.Use(tracing.NewPlugin(tracing.WithDBName("ppt")))
	//if err != nil {
	//	return nil, err
	//}

	if err != nil {
		return nil, err
	}
	return &ppt{
		db: db,
	}, nil
}

func (p ppt) Save(ctx context.Context, node *entity.Node) apierrors.ApiError {
	result := p.db.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(contracts.MapFromEntity(node))
	if result.Error != nil {
		return apierrors.NewInternalServerApiError(result.Error.Error())
	}
	return nil
}

func (p ppt) GetById(ctx context.Context, id string) (*entity.Node, apierrors.ApiError) {
	node := &contracts.Node{ID: id}
	result := p.db.WithContext(ctx).Clauses().Preload("Children", "parent_id = ? and id <> ?", id, id).First(node)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewNotFoundApiError("node not found")
		}
		return nil, apierrors.NewInternalServerApiError(result.Error.Error())
	}
	return contracts.MapToEntity(node), nil
}

func (p ppt) GetTree(ctx context.Context, rootId string) (*entity.Node, apierrors.ApiError) {
	sql := `
		WITH RECURSIVE node_tree AS (
			SELECT id, parent_id, 0 as level
			FROM nodes
			WHERE id = ?
		  
			UNION ALL
		  
			SELECT n.id, n.parent_id, nt.level + 1 as level
			FROM nodes n
			INNER JOIN node_tree nt ON n.parent_id = nt.id
			WHERE n.id <> n.parent_id
		)
		SELECT * FROM node_tree;
	`
	rows, err := p.db.WithContext(ctx).Raw(sql, rootId).Rows()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewNotFoundApiError("node not found")
		}
		return nil, apierrors.NewInternalServerApiError(err.Error())
	}
	defer rows.Close()

	var nodes []*contracts.Node
	for rows.Next() {
		var node contracts.Node
		err = p.db.WithContext(ctx).ScanRows(rows, &node)
		if err != nil {
			return nil, apierrors.NewInternalServerApiError(err.Error())
		}
		nodes = append(nodes, &node)
	}

	nodeMap := make(map[string]*entity.Node)
	childMap := make(map[string][]*entity.Node)
	var root *entity.Node

	for _, node := range nodes {
		entityNode := contracts.MapToEntity(node)
		nodeMap[entityNode.Id] = entityNode
		childMap[entityNode.ParentId] = append(childMap[entityNode.ParentId], entityNode)
	}

	for id, node := range nodeMap {
		node.Children = childMap[id]

		if node.Id == rootId {
			root = node
		}
	}

	if root == nil {
		return nil, apierrors.NewNotFoundApiError("root node not found")
	}

	return root, nil
}
