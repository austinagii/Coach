package api

import ()

type ErrorCode string

const (
	InternalError  ErrorCode = "internal_error"
	InvalidRequest ErrorCode = "invalid_request"
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
