package nested

import (
	"github.com/F-Amaral/tcc/constants"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	fx.Annotate(NewNestedHandler, fx.ParamTags(constants.NestedServiceName)),
)
