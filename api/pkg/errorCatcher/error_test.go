package errorCatcher

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ApiErrorSuite struct {
	suite.Suite
}

// Test ApiError.Error() with cause
func (suite *ApiErrorSuite) TestApiError_Error_WithCause() {
	cause := errors.New("underlying error")
	err := &ApiError{
		Code:    "TEST_ERROR",
		Message: "test message",
		Cause:   cause,
	}

	expected := "TEST_ERROR: test message: underlying error"
	suite.Equal(expected, err.Error())
}

// Test ApiError.Error() without cause
func (suite *ApiErrorSuite) TestApiError_Error_WithoutCause() {
	err := &ApiError{
		Code:    "TEST_ERROR",
		Message: "test message",
		Cause:   nil,
	}

	expected := "TEST_ERROR: test message"
	suite.Equal(expected, err.Error())
}

// Test ApiError.Unwrap()
func (suite *ApiErrorSuite) TestApiError_Unwrap() {
	cause := errors.New("underlying error")
	err := &ApiError{
		Code:  "TEST_ERROR",
		Cause: cause,
	}

	suite.Equal(cause, err.Unwrap())
}

// Test ApiError.Unwrap() with nil cause
func (suite *ApiErrorSuite) TestApiError_Unwrap_NilCause() {
	err := &ApiError{
		Code:  "TEST_ERROR",
		Cause: nil,
	}

	suite.Nil(err.Unwrap())
}

// Test NewAuthError
func (suite *ApiErrorSuite) TestNewAuthError() {
	cause := errors.New("invalid token")
	err := NewAuthError("authentication failed", cause)

	suite.Equal("AUTH_ERROR", err.Code)
	suite.Equal("authentication failed", err.Message)
	suite.Equal(http.StatusUnauthorized, err.StatusCode)
	suite.Equal(cause, err.Cause)
}

// Test NewAuthError without cause
func (suite *ApiErrorSuite) TestNewAuthError_NoCause() {
	err := NewAuthError("authentication failed", nil)

	suite.Equal("AUTH_ERROR", err.Code)
	suite.Equal("authentication failed", err.Message)
	suite.Equal(http.StatusUnauthorized, err.StatusCode)
	suite.Nil(err.Cause)
}

// Test NewPermissionError
func (suite *ApiErrorSuite) TestNewPermissionError() {
	cause := errors.New("insufficient permissions")
	err := NewPermissionError("access denied", cause)

	suite.Equal("PERMISSION_DENIED", err.Code)
	suite.Equal("access denied", err.Message)
	suite.Equal(http.StatusForbidden, err.StatusCode)
	suite.Equal(cause, err.Cause)
}

// Test NewValidationError
func (suite *ApiErrorSuite) TestNewValidationError() {
	cause := errors.New("invalid email format")
	err := NewValidationError("validation failed", cause)

	suite.Equal("VALIDATION_ERROR", err.Code)
	suite.Equal("validation failed", err.Message)
	suite.Equal(http.StatusBadRequest, err.StatusCode)
	suite.Equal(cause, err.Cause)
}

// Test NewNotFoundError
func (suite *ApiErrorSuite) TestNewNotFoundError() {
	cause := errors.New("user not found")
	err := NewNotFoundError("resource not found", cause)

	suite.Equal("NOT_FOUND", err.Code)
	suite.Equal("resource not found", err.Message)
	suite.Equal(http.StatusNotFound, err.StatusCode)
	suite.Equal(cause, err.Cause)
}

// Test NewDatabaseError
func (suite *ApiErrorSuite) TestNewDatabaseError() {
	cause := errors.New("unique constraint violation")
	err := NewDatabaseError("database operation failed", cause)

	suite.Equal("DATABASE_ERROR", err.Code)
	suite.Equal("database operation failed", err.Message)
	suite.Equal(http.StatusUnprocessableEntity, err.StatusCode)
	suite.Equal(cause, err.Cause)
}

// Test NewInternalError
func (suite *ApiErrorSuite) TestNewInternalError() {
	cause := errors.New("unexpected panic")
	err := NewInternalError("internal server error", cause)

	suite.Equal("INTERNAL_ERROR", err.Code)
	suite.Equal("internal server error", err.Message)
	suite.Equal(http.StatusInternalServerError, err.StatusCode)
	suite.Equal(cause, err.Cause)
}

// Test errors.Is() compatibility
func (suite *ApiErrorSuite) TestApiError_ErrorsIs() {
	baseErr := errors.New("base error")
	apiErr := NewAuthError("auth failed", baseErr)

	suite.True(errors.Is(apiErr, baseErr))
}

// Test errors.As() compatibility
func (suite *ApiErrorSuite) TestApiError_ErrorsAs() {
	err := NewAuthError("auth failed", nil)

	var apiErr *ApiError
	suite.True(errors.As(err, &apiErr))
	suite.Equal("AUTH_ERROR", apiErr.Code)
}

func TestApiErrorSuite(t *testing.T) {
	suite.Run(t, new(ApiErrorSuite))
}
