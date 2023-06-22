package middlewares

import (
	"github.com/F-Amaral/tcc/internal/log"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
)

func LogInContextMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.LogCtxKey, logger)
		c.Next()
	}
}
