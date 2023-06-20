package service

import (
	"github.com/F-Amaral/tcc/constants"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(NewPpt, fx.ParamTags(constants.PPTRepositoryName), fx.ResultTags(constants.PPTServiceName)),
	fx.Annotate(NewNested, fx.ParamTags(constants.NestedRepositoryName), fx.ResultTags(constants.NestedServiceName)),
)
