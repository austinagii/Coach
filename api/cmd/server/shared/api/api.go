package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type ErrorCode string

const (
	ErrorCodeGeneralError ErrorCode = "gen_server_Error"
	ErrorCodeBadRequest   ErrorCode = "bad_request"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewApiError(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:    string(code),
		Message: message,
	}
}

func BindApiErrorResponse(
	c *gin.Context,
	errorMessage string,
	statusCode int,
	errorCode ErrorCode,
	err error,
) {
	slog.Error(errorMessage, "error", err)
	c.IndentedJSON(statusCode, NewApiError(
		errorCode,
		fmt.Errorf("%s: %w", errorMessage, err).Error(),
	))
}
