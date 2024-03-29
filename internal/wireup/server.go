package wireup

import (
	"context"
	"errors"
	"fmt"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/F-Amaral/tcc/internal/wireup/middlewares"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"go.uber.org/fx"
	"time"
)

const (
	defaultShutdownDelay = 30 * time.Second
)

var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Invoke(InitServer, log.NewLogger),
)

type Server struct {
	Engine *gin.Engine
	Tracer telemetry.Telemetry
}

func NewServer(logger log.Logger, telemetry telemetry.Telemetry) *Server {
	engine := gin.New()
	engine.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	engine.Use(nrgin.Middleware(telemetry))
	engine.Use(middlewares.TracerInContextMiddleware(telemetry))
	engine.Use(middlewares.LogInContextMiddleware(logger))
	server := Server{
		Engine: engine,
		Tracer: telemetry,
	}
	return &server
}

func InitServer(server *Server, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go runServer(server)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func runServer(server *Server) {
	err := server.Engine.Run(":8080")
	if err != nil {
		ScheduleShutdown(err.Error(), err)
	}
}

func ScheduleShutdown(reason string, cause error) {
	log.Error(context.Background()).LogError(errors.New(fmt.Sprintf("server failed to start, scheduling shutdown in %s for reason %s", defaultShutdownDelay, reason)))
	time.Sleep(defaultShutdownDelay)
	panic(cause)
}
