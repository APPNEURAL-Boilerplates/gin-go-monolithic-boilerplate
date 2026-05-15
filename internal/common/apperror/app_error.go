package apperror

import (
	stderrors "errors"
	"net/http"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Details any
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(status int, code string, message string, details any) *AppError {
	return &AppError{Status: status, Code: code, Message: message, Details: details}
}

func BadRequest(message string, details any) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Validation(message string, details any) *AppError {
	return New(http.StatusBadRequest, "VALIDATION_ERROR", message, details)
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message, nil)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, nil)
}

func Internal() *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred.", nil)
}

func From(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr
	}

	return Internal()
}
