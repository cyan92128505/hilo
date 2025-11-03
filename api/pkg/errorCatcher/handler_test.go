package errorCatcher

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type ErrorHandlerMiddlewareSuite struct {
	suite.Suite
	obLog  *observer.ObservedLogs
	logger *zap.Logger
}

func (suite *ErrorHandlerMiddlewareSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	suite.logger = zap.New(observedZapCore)
	suite.obLog = observedLogs
}

// Test ErrorHandlerMiddleware with ApiError - Auth Error
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_ApiError_Auth() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.Error(NewAuthError("invalid token", errors.New("jwt expired")))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusUnauthorized, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("request failed", firstLog.Message)
	suite.Equal("AUTH_ERROR", firstLog.ContextMap()["code"])
	suite.Equal("invalid token", firstLog.ContextMap()["message"])
}

// Test ErrorHandlerMiddleware with ApiError - Validation Error
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_ApiError_Validation() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.POST("/test", func(c *gin.Context) {
		c.Error(NewValidationError("invalid email", errors.New("must be valid email format")))
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusBadRequest, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("request failed", firstLog.Message)
	suite.Equal("VALIDATION_ERROR", firstLog.ContextMap()["code"])
}

// Test ErrorHandlerMiddleware with ApiError - Permission Error
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_ApiError_Permission() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.Error(NewPermissionError("access denied", nil))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusForbidden, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("request failed", firstLog.Message)
	suite.Equal("PERMISSION_DENIED", firstLog.ContextMap()["code"])
}

// Test ErrorHandlerMiddleware with ApiError - Not Found Error
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_ApiError_NotFound() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.Error(NewNotFoundError("user not found", nil))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusNotFound, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("request failed", firstLog.Message)
	suite.Equal("NOT_FOUND", firstLog.ContextMap()["code"])
}

// Test ErrorHandlerMiddleware with ApiError - Database Error
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_ApiError_Database() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.POST("/test", func(c *gin.Context) {
		c.Error(NewDatabaseError("unique constraint violation", errors.New("duplicate key")))
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusUnprocessableEntity, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("request failed", firstLog.Message)
	suite.Equal("DATABASE_ERROR", firstLog.ContextMap()["code"])
}

// Test ErrorHandlerMiddleware with ApiError - Internal Error
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_ApiError_Internal() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.Error(NewInternalError("unexpected error", errors.New("nil pointer")))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusInternalServerError, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("request failed", firstLog.Message)
	suite.Equal("INTERNAL_ERROR", firstLog.ContextMap()["code"])
}

// Test ErrorHandlerMiddleware with unknown error (not ApiError)
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_UnknownError() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("some random error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusInternalServerError, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("unexpected error", firstLog.Message)
}

// Test ErrorHandlerMiddleware with no errors
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_NoError() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(http.StatusOK, result.StatusCode)
	suite.Equal(0, suite.obLog.Len())
}

// Test ErrorHandlerMiddleware logs request path and method
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_LogsPathAndMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.POST("/api/users", func(c *gin.Context) {
		c.Error(NewAuthError("auth failed", nil))
	})

	req := httptest.NewRequest(http.MethodPost, "/api/users", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("/api/users", firstLog.ContextMap()["path"])
	suite.Equal("POST", firstLog.ContextMap()["method"])
}

// Test ErrorHandlerMiddleware with multiple errors (uses last one)
func (suite *ErrorHandlerMiddlewareSuite) TestErrorHandlerMiddleware_MultipleErrors_UsesLast() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(ErrorHandlerMiddleware(suite.logger))
	route.GET("/test", func(c *gin.Context) {
		c.Error(NewAuthError("first error", nil))
		c.Error(NewValidationError("second error", nil))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()

	// Should use the last error (ValidationError -> 400)
	suite.Equal(http.StatusBadRequest, result.StatusCode)
	suite.Equal(1, suite.obLog.Len())

	firstLog := suite.obLog.All()[0]
	suite.Equal("VALIDATION_ERROR", firstLog.ContextMap()["code"])
	suite.Equal("second error", firstLog.ContextMap()["message"])
}

func TestErrorHandlerMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(ErrorHandlerMiddlewareSuite))
}
