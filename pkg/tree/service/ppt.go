package service

import (
	"context"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
	"github.com/F-Amaral/tcc/pkg/tree/domain/repositories"
	"github.com/F-Amaral/tcc/pkg/tree/domain/services"
	"net/http"
)

type ppt struct {
	repository repositories.PPTTree
}

func NewPpt(repository repositories.PPTTree) services.Tree {
	return &ppt{
		repository: repository,
	}
}

func (p ppt) Create(ctx context.Context, id string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Service Create")
	defer tx.End()
	node := &entity.Node{
		Id:       id,
		ParentId: id,
	}

	err := p.repository.Save(ctx, node)
	if err != nil {
		log.WrapTransaction(tx).Error(err)
		return nil, err
	}
	return node, nil
}

func (p ppt) GetTree(ctx context.Context, nodeId string, recursive bool) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Service GetTreeRecursive")
	defer tx.End()
	if recursive {
		node, err := p.repository.GetTreeRecursive(ctx, nodeId)
		if err != nil {
			log.WrapTransaction(tx).Error(err)
			return nil, err
		}
		return node, nil
	}
	node, err := p.repository.GetTree(ctx, nodeId)
	if err != nil {
		log.WrapTransaction(tx).Error(err)
		return nil, err
	}
	return node, nil
}

func (p ppt) AddToParent(ctx context.Context, parentId, childId string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Service AddToParent")
	defer tx.End()
	childNode, err := p.getOrCreate(ctx, childId)
	if err != nil {
		log.WrapTransaction(tx).Error(err)
		return nil, err
	}

	if childNode.ParentId != "" && childNode.ParentId != childNode.Id {
		if childNode.ParentId != parentId {
			return nil, apierrors.NewBadRequestError("node already has a parent")
		}
		return nil, apierrors.NewBadRequestError("node already has this parent")
	}

	parentNode, err := p.getOrCreate(ctx, parentId)
	if err != nil {
		return nil, err
	}

	childNode.ParentId = parentId
	saveErr := p.repository.Save(ctx, childNode)
	if saveErr != nil {
		log.WrapTransaction(tx).Error(err)
		return nil, err
	}

	parentNode.Children = append(parentNode.Children, childNode)
	return parentNode, nil
}

func (p ppt) RemoveFromParent(ctx context.Context, parentId string, childId string) (*entity.Node, apierrors.ApiError) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Service RemoveFromParent")
	defer tx.End()
	parentNode, err := p.repository.GetById(ctx, parentId)
	if err != nil {
		return nil, err
	}

	childNode, err := p.repository.GetById(ctx, childId)
	if err != nil {
		return nil, err
	}

	if childNode.ParentId != parentId {
		return nil, apierrors.NewBadRequestError("node does not have this parent")
	}

	childNode.ParentId = childNode.Id
	saveErr := p.repository.Save(ctx, childNode)
	if saveErr != nil {
		return nil, err
	}

	for i, child := range parentNode.Children {
		if child.Id == childId {
			parentNode.Children = append(parentNode.Children[:i], parentNode.Children[i+1:]...)
			break
		}
	}

	return parentNode, nil
}

func (p ppt) getOrCreate(ctx context.Context, id string) (*entity.Node, apierrors.ApiError) {
	node, err := p.repository.GetById(ctx, id)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			return nil, err
		}
		return p.Create(ctx, id)
	}

	return node, nil
}
