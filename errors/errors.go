package errors

import (
	"net/http"
)

type CustomError struct {
	HTTPStatus  int    `json:"-"`
	Code        string `json:"code"`
	Message     string `json:"message"`
	IsRetryable bool   `json:"-"`
}

func (a *CustomError) Error() string {
	return a.Message
}

func (customError *CustomError) Retryable() {
	customError.IsRetryable = true
}

func (customError *CustomError) NotRetryable() {
	customError.IsRetryable = false
}

func NewCustomError(httpStatus int, code string, err error) *CustomError {
	return &CustomError{
		HTTPStatus:  httpStatus,
		Code:        code,
		Message:     err.Error(),
		IsRetryable: true,
	}
}

func NewUnauthorizedError(err error) *CustomError {
	return &CustomError{
		HTTPStatus:  http.StatusUnauthorized,
		Code:        "UNAUTHORIZED",
		Message:     err.Error(),
		IsRetryable: false,
	}
}

func NewForbiddenError(err error) *CustomError {
	return &CustomError{
		HTTPStatus:  http.StatusForbidden,
		Code:        "FORBIDDEN",
		Message:     err.Error(),
		IsRetryable: false,
	}
}

func NewInternalServerError(code string, err error) *CustomError {
	return &CustomError{
		HTTPStatus:  http.StatusInternalServerError,
		Code:        code,
		Message:     err.Error(),
		IsRetryable: false,
	}
}

func NewNotFoundError(code string, err error) *CustomError {
	return &CustomError{
		HTTPStatus:  http.StatusNotFound,
		Code:        code,
		Message:     err.Error(),
		IsRetryable: true,
	}
}

func NewBadRequest(code string, err error) *CustomError {
	return &CustomError{
		HTTPStatus:  http.StatusBadRequest,
		Code:        code,
		Message:     err.Error(),
		IsRetryable: false,
	}
}

func NewInvalidPayloadError(code string, err error) *CustomError {
	return &CustomError{
		HTTPStatus:  http.StatusBadRequest,
		Code:        code,
		Message:     err.Error(),
		IsRetryable: false,
	}
}

func NewUnkownDatabaseError(err error) *CustomError {
	return NewInternalServerError("UNKNOWN_DATABASE_ERROR", err)
}
