package ppt

import (
	"github.com/F-Amaral/tcc/internal/wireup"
)

func RegisterHandlers(server *wireup.Server, handler *PptHandler) {
	server.Engine.GET("/ppt/:id", handler.GetTree)
	server.Engine.POST("/ppt/:parentId/:childId", handler.AddToParent)
	server.Engine.DELETE("/ppt/:parentId/:childId", handler.RemoveFromParent)

}
