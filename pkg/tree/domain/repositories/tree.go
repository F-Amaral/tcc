package repositories

import (
	"context"
	"github.com/F-Amaral/tcc/internal/apierrors"
	"github.com/F-Amaral/tcc/pkg/tree/domain/entity"
)

type Tree interface {
	Save(context.Context, *entity.Node) apierrors.ApiError
	GetById(context.Context, string) (*entity.Node, apierrors.ApiError)
	GetTree(context.Context, string) (*entity.Node, apierrors.ApiError)
}
