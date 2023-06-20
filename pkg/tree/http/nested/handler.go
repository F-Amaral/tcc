package nested

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/F-Amaral/tcc/pkg/tree/domain/services"
	"github.com/F-Amaral/tcc/pkg/tree/http/nested/contracts"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"io"
	"mime/multipart"
)

func NewNestedHandler(nested services.Tree) *NestedHandler {
	return &NestedHandler{nested: nested}
}

type NestedHandler struct {
	nested services.Tree
}

func (handler *NestedHandler) GetTree(ctx *gin.Context) {
	request := contracts.GetTreeRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Id == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing id"))
		return
	}

	res, err := handler.nested.GetTree(ctx, request.Id)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (handler *NestedHandler) AddToParent(ctx *gin.Context) {
	request := contracts.AddToParentRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.ParentId == "" || request.NodeId == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing parent or child id"))
		return
	}

	res, err := handler.nested.AddToParent(ctx, request.ParentId, request.NodeId)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (handler *NestedHandler) RemoveFromParent(ctx *gin.Context) {
	request := contracts.RemoveFromParentRequest{}
	if err := ctx.BindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.ParentId == "" || request.NodeId == "" {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("missing parent or child id"))
		return
	}

	res, err := handler.nested.RemoveFromParent(ctx, request.ParentId, request.NodeId)
	if err != nil {
		ctx.JSON(err.Status(), err)
		ctx.Writer.WriteHeaderNow()
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (handler *NestedHandler) UploadCSV(ctx *gin.Context) {
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

func (handler *NestedHandler) processCSV(ctx context.Context, file multipart.File) error {
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

		_, err = handler.nested.AddToParent(ctx, parentId, nodeId)
		if err != nil {
			return err
		}
	}

	return nil
}
