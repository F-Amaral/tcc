package ppt

import (
	"context"
	"database/sql"
	"errors"
	"github.com/F-Amaral/tcc/constants"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/pkg/tree/domain/repositories"
	"github.com/F-Amaral/tcc/pkg/tree/repository/mysql/ppt/contracts"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/mysql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	logger2 "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

type ppt struct {
	db     *gorm.DB
	tracer *newrelic.Application
}

func NewPpt(config *viper.Viper, logger log.Logger, nr *newrelic.Application) (repositories.Tree, error) {
	logWrap := zapgorm2.New(logger.Desugar())
	logWrap.LogMode(logger2.Silent)
	logWrap.SetAsDefault()
	nrDb, err := sql.Open("nrmysql", config.GetString(constants.PPtDbDsnKey))
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: nrDb,
	}), &gorm.Config{Logger: logWrap})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&contracts.NodeParent{}, &contracts.Node{})

	if err != nil {
		return nil, err
	}
	return &ppt{
		db:     db,
		tracer: nr,
	}, nil
}

func (p ppt) Save(ctx context.Context, node *entity.Node) apierrors.ApiError {
	trace := p.tracer.StartTransaction("PPT Save")
	traceCtx := newrelic.NewContext(ctx, trace)
	defer trace.End()
	result := p.db.WithContext(traceCtx).Clauses(clause.OnConflict{UpdateAll: true}).Create(contracts.MapFromEntity(node))
	if result.Error != nil {
		return apierrors.NewInternalServerApiError(result.Error.Error())
	}
	return nil
}

func (p ppt) GetById(ctx context.Context, id string) (*entity.Node, apierrors.ApiError) {
	trace := p.tracer.StartTransaction("PPT GetById")
	traceCtx := newrelic.NewContext(ctx, trace)
	defer trace.End()
	node := &contracts.NodeParent{ID: id}
	result := p.db.WithContext(traceCtx).Clauses().Preload("Children", "parent_id = ? and id <> ?", id, id).First(node)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewNotFoundApiError("node not found")
		}
		return nil, apierrors.NewInternalServerApiError(result.Error.Error())
	}
	return contracts.MapToEntity(node), nil
}

func (p ppt) GetTreeRecursive(ctx context.Context, rootId string) (*entity.Node, apierrors.ApiError) {
	sql := `
		WITH RECURSIVE node_tree AS (
			SELECT id, parent_id, 0 as level
			FROM nodes as n
			left join container_parent cp on cp.id = n.id
			WHERE id = ?
			UNION ALL
			SELECT n.id, n.parent_id, nt.level + 1 as level
			FROM nodes n
			INNER JOIN node_tree nt ON n.parent_id = nt.id
			WHERE n.id <> n.parent_id
		)
		SELECT * FROM node_tree;
	`
	trace := p.tracer.StartTransaction("Ppt GetTreeRecursive")
	traceCtx := newrelic.NewContext(ctx, trace)
	defer trace.End()
	rows, err := p.db.WithContext(traceCtx).Raw(sql, rootId).Rows()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewNotFoundApiError("node not found")
		}
		return nil, apierrors.NewInternalServerApiError(err.Error())
	}
	defer rows.Close()

	buildTrace := p.tracer.StartTransaction("Ppt GetTreeRecursive Build")
	buildTraceCtx := newrelic.NewContext(traceCtx, buildTrace)
	defer buildTrace.End()
	var nodes []*contracts.NodeParent
	for rows.Next() {
		var node contracts.NodeParent
		err = p.db.WithContext(buildTraceCtx).ScanRows(rows, &node)
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
