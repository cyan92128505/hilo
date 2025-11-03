package errorCatcher

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	Code       string
	Message    string
	StatusCode int
	Cause      error
}

func (e *ApiError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ApiError) Unwrap() error {
	return e.Cause
}

// Error constructors
func NewAuthError(message string, cause error) *ApiError {
	return &ApiError{
		Code:       "AUTH_ERROR",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Cause:      cause,
	}
}

func NewPermissionError(message string, cause error) *ApiError {
	return &ApiError{
		Code:       "PERMISSION_DENIED",
		Message:    message,
		StatusCode: http.StatusForbidden,
		Cause:      cause,
	}
}

func NewValidationError(message string, cause error) *ApiError {
	return &ApiError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}

func NewNotFoundError(message string, cause error) *ApiError {
	return &ApiError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
		Cause:      cause,
	}
}

func NewDatabaseError(message string, cause error) *ApiError {
	return &ApiError{
		Code:       "DATABASE_ERROR",
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
		Cause:      cause,
	}
}

func NewInternalError(message string, cause error) *ApiError {
	return &ApiError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Cause:      cause,
	}
}
