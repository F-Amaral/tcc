package services

import (
	"context"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
)

type Tree interface {
	Create(context.Context, string) (*entity.Node, apierrors.ApiError)
	GetTree(context.Context, string) (*entity.Node, apierrors.ApiError)
	AddToParent(context.Context, string, string) (*entity.Node, apierrors.ApiError)
	RemoveFromParent(context.Context, string, string) (*entity.Node, apierrors.ApiError)
}
