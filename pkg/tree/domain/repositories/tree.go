package repositories

import (
	"context"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
)

type Tree interface {
	Save(context.Context, *entity.Node) apierrors.ApiError
	GetById(context.Context, string) (*entity.Node, apierrors.ApiError)
	GetTreeRecursive(context.Context, string) (*entity.Node, apierrors.ApiError)
}

type NestedTree interface {
	Tree
	AppendToTree(context.Context, string, *entity.Node) (*entity.Node, apierrors.ApiError)
}
