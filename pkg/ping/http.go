package ping

import (
	"github.com/F-Amaral/tcc/internal/log"
	"github.com/F-Amaral/tcc/internal/web"
	"github.com/F-Amaral/tcc/internal/wireup"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"go.uber.org/fx"
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
