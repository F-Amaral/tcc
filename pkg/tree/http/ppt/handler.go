package ppt

import (
	"errors"
	"github.com/F-Amaral/tcc/pkg/tree/domain/services"
	"github.com/F-Amaral/tcc/pkg/tree/http/ppt/contracts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewPptHandler(ppt services.Tree) *PptHandler {
	return &PptHandler{ppt: ppt}
}

type PptHandler struct {
	ppt services.Tree
}

func (handler *PptHandler) GetTree(ctx *gin.Context) {
	request := contracts.GetTreeRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Id == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing id"))
		return
	}

	res, err := handler.ppt.GetTree(ctx, request.Id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (handler *PptHandler) AddToParent(ctx *gin.Context) {
	request := contracts.AddToParentRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.ParentId == "" || request.NodeId == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing parent or child id"))
		return
	}

	res, err := handler.ppt.AddToParent(ctx, request.ParentId, request.NodeId)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (handler *PptHandler) RemoveFromParent(ctx *gin.Context) {
	request := contracts.RemoveFromParentRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.ParentId == "" || request.NodeId == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing parent or child id"))
		return
	}

	res, err := handler.ppt.RemoveFromParent(ctx, request.ParentId, request.NodeId)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}
