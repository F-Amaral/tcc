package wireup

import (
	"context"
	"errors"
	"fmt"
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"time"
)

const (
	defaultShutdownDelay = 30 * time.Second
)

var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Invoke(InitServer),
)

type Server struct {
	Engine *gin.Engine
}

func NewServer() *Server {
	server := Server{Engine: gin.Default()}
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
	log.Err(errors.New(fmt.Sprintf("server failed to start, scheduling shutdown in %s for reason %s", defaultShutdownDelay, reason)))
	time.Sleep(defaultShutdownDelay)
	panic(cause)
}
