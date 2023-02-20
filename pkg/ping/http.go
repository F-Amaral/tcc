package ping

import (
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/web"
	"github.com/F-Amaral/tcc/internal/wireup"
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
