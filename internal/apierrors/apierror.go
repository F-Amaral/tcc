package apierrors

import (
	"context"
	"fmt"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
)

type ApiError interface {
	Status() int
	Code() string
	Msg() string
	Error() string
}

type apiError struct {
	Message   string `json:"msg"`
	ErrCode   string `json:"code"`
	ErrStatus int    `json:"status"`
}

func BuildErrorWithContext(ctx context.Context, err error) ApiError {
	return &apiError{
		Message:   err.Error(),
		ErrCode:   err.Error(),
		ErrStatus: http.StatusInternalServerError,
	}
}

func NewInternalServerApiError(msg string) ApiError {
	return &apiError{
		Message:   msg,
		ErrCode:   "internal_server_error",
		ErrStatus: http.StatusInternalServerError,
	}
}

func NewBadRequestError(msg string) ApiError {
	return &apiError{
		Message:   msg,
		ErrCode:   "bad_request",
		ErrStatus: http.StatusBadRequest,
	}
}

func NewNotFoundApiError(msg string) ApiError {
	return &apiError{
		Message:   msg,
		ErrCode:   "not_found",
		ErrStatus: http.StatusNotFound,
	}
}

func (a apiError) Status() int {
	return a.ErrStatus
}

func (a apiError) Code() string {
	return a.ErrCode
}

func (a apiError) Msg() string {
	return a.Message
}

func (a apiError) Error() string {
	return fmt.Sprintf("[%d] - Code: %s, Message: %s", a.ErrStatus, a.ErrCode, a.Msg)
}
