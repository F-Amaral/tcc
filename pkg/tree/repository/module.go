package repository

import (
	"github.com/F-Amaral/tcc/pkg/tree/repository/mysql/nested"
	"github.com/F-Amaral/tcc/pkg/tree/repository/mysql/ppt"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(ppt.NewPpt),
	fx.Annotate(nested.NewNested),
)
