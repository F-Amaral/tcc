package ping

import (
	"git.jetbrains.space/philldev/tcc/internal/log"
	"git.jetbrains.space/philldev/tcc/internal/web"
	"git.jetbrains.space/philldev/tcc/internal/wireup"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Invoke(DefinePingEndpoint)

const (
	baseUrl = "/ping"
)

func DefinePingEndpoint(server *wireup.Server) {
	server.Engine.GET(baseUrl, gin.WrapF(MakePingHandler))
}

func MakePingHandler(w http.ResponseWriter, r *http.Request) {
	err := web.EncodeJson(w, "pong", http.StatusOK)
	if err != nil {
		log.Err(err)
	}
}
