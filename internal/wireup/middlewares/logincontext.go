package middlewares

import (
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/gin-gonic/gin"
)

func LogInContextMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(log.LogCtxKey, logger)
		c.Next()
	}
}

func TracerInContextMiddleware(t telemetry.Telemetry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(telemetry.TelemetryCtxKey, t)
		c.Next()
	}
}
