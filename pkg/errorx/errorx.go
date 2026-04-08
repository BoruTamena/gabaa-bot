package errorx

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrForbidden       ErrorCode = "FORBIDDEN"
	ErrValidation      ErrorCode = "VALIDATION_ERROR"
	ErrNotFound        ErrorCode = "NOT_FOUND"
	ErrInternal        ErrorCode = "INTERNAL_ERROR"
	ErrBadRequest      ErrorCode = "BAD_REQUEST"
	ErrConflict        ErrorCode = "CONFLICT"
)

type AppError struct {
	Code    ErrorCode              `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Status  int                    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code ErrorCode, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message
	return e
}

func ErrorResponse(err error) (int, *AppError) {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status, appErr
	}

	return http.StatusInternalServerError, &AppError{
		Code:    ErrInternal,
		Message: "An internal server error occurred",
		Status:  http.StatusInternalServerError,
	}
}
