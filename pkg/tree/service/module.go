package service

import (
	"github.com/F-Amaral/tcc/constants"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(NewPpt, fx.ResultTags(constants.PPTServiceName)),
	fx.Annotate(NewNested, fx.ResultTags(constants.NestedServiceName)),
)
