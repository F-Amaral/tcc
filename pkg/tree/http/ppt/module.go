package ppt

import (
	"github.com/F-Amaral/tcc/constants"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(NewPptHandler, fx.ParamTags(constants.PPTServiceName)),
)
