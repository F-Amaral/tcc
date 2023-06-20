package repository

import (
	"github.com/F-Amaral/tcc/constants"
	"github.com/F-Amaral/tcc/pkg/tree/repository/mysql"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(mysql.NewPpt, fx.ResultTags(constants.PPTRepositoryName)),
	fx.Annotate(mysql.NewNested, fx.ResultTags(constants.NestedRepositoryName)),
)
