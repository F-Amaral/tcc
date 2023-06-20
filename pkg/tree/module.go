package tree

import (
	"github.com/F-Amaral/tcc/pkg/tree/http/ppt"
	"go.uber.org/fx"
)

var Module = fx.Invoke(ppt.RegisterHandlers)
