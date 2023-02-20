package apierrors

import (
	"context"
	"github.com/F-Amaral/tcc/internal/log"
	"net/http"
)

type ApiError struct {
	Msg    string `json:"msg"`
	Code   string `json:"code"`
	Status int    `json:"status"`
}

func BuildErrorWithContext(ctx context.Context, err error) *ApiError {
	log.Err(err)
	return &ApiError{
		Msg:    err.Error(),
		Code:   err.Error(),
		Status: http.StatusInternalServerError,
	}
}
