package nested

import (
	"github.com/F-Amaral/tcc/internal/wireup"
)

func RegisterHandlers(server *wireup.Server, handler *NestedHandler) {
	server.Engine.GET("/nested/:id", handler.GetTree)
	server.Engine.POST("/nested/:parentId/:childId", handler.AddToParent)
	server.Engine.POST("/nested/upload", handler.UploadCSV)
	server.Engine.DELETE("/nested/:parentId/:childId", handler.RemoveFromParent)

}
