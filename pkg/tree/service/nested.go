package service

import (
	"context"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/pkg/tree/domain/repositories"
	"github.com/F-Amaral/tcc/pkg/tree/domain/services"
	"net/http"
)

type nested struct {
	repository repositories.NestedTree
}

func NewNested(repository repositories.NestedTree) services.Tree {
	return &nested{
		repository: repository,
	}
}

func (p nested) Create(ctx context.Context, id string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Nested Service Create")
	defer tx.End()
	node := &entity.Node{
		Id: id,
	}

	err := p.repository.Save(ctx, node)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p nested) GetTree(ctx context.Context, nodeId string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Nested Service GetTree")
	defer tx.End()
	return p.repository.GetTree(ctx, nodeId)
}

func (p nested) AddToParent(ctx context.Context, parentId, childId string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Nested Service AddToParent")
	defer tx.End()
	parentNode, err := p.getOrCreate(ctx, parentId)
	if err != nil {
		return nil, err
	}

	childNode, err := p.getOrCreate(ctx, childId)
	if err != nil {
		return nil, err
	}

	if childNode.ParentId != "" {
		if childNode.ParentId != parentId {
			return nil, apierrors.NewBadRequestError("node already has a parent")
		}
		return nil, apierrors.NewBadRequestError("node already has this parent")
	}

	parentNode, err = p.repository.AppendToTree(ctx, parentNode.Id, childNode)
	if err != nil {
		return nil, err
	}

	return parentNode, nil
}

func (p nested) RemoveFromParent(ctx context.Context, _, nodeId string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Nested Service RemoveFromParent")
	defer tx.End()
	node, err := p.repository.GetById(ctx, nodeId)
	if err != nil {
		return nil, err
	}

	node.ParentId = ""
	saveErr := p.repository.Save(ctx, node)
	if saveErr != nil {
		return nil, err
	}

	return node, nil
}

func (p nested) getOrCreate(ctx context.Context, id string) (*entity.Node, apierrors.ApiError) {
	node, err := p.repository.GetById(ctx, id)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
		return p.Create(ctx, id)
	}

	return node, nil
}
