package ppt

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/F-Amaral/tcc/internal/telemetry"
	"github.com/F-Amaral/tcc/pkg/tree/domain/services"
	"github.com/F-Amaral/tcc/pkg/tree/http/ppt/contracts"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"io"
	"mime/multipart"
)

func NewPptHandler(ppt services.Tree) *PptHandler {
	return &PptHandler{ppt: ppt}
}

type PptHandler struct {
	ppt services.Tree
}

func (handler *PptHandler) GetTree(ctx *gin.Context) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Handler GetTreeRecursive")
	defer tx.End()
	request := contracts.GetTreeRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Id == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing id"))
		return
	}

	recursive := true
	if recursiveQuery := ctx.Query("recursive"); recursiveQuery != "" {
		r, err := strconv.ParseBool(recursiveQuery)
		if err == nil {
			recursive = r
		}
	}

	res, err := handler.ppt.GetTree(ctx, request.Id, recursive)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (handler *PptHandler) AddToParent(ctx *gin.Context) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Handler AddToParent")
	defer tx.End()
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
	tx := telemetry.With(ctx).StartTransaction("Ppt Handler RemoveFromParent")
	defer tx.End()
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

func (handler *PptHandler) UploadCSV(ctx *gin.Context) {
	tx := telemetry.With(ctx).StartTransaction("Ppt Handler UploadCSV")
	defer tx.End()
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	err = handler.processCSV(ctx, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (handler *PptHandler) processCSV(ctx context.Context, file multipart.File) error {
	reader := csv.NewReader(file)

	// Skip the header
	_, err := reader.Read()
	if err != nil && err != io.EOF {
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		nodeId := record[0]
		parentId := record[1]

		_, err = handler.ppt.AddToParent(ctx, parentId, nodeId)
		if err != nil {
			return err
		}
	}

	return nil
}
