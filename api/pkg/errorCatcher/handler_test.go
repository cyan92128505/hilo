package errorCatcher

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type HandlerSuite struct {
	suite.Suite
	obLog  *observer.ObservedLogs
	logger *zap.Logger
}

func (suite *HandlerSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	suite.logger = zap.New(observedZapCore)
	suite.obLog = observedLogs
}

func (suite *HandlerSuite) TestPanicErrorHandler_PanicNormalError_ShouldMatchExpected() {
	func() {
		defer PanicErrorHandler(suite.logger, "Mock test: ")
		panic(errors.New("got error"))
	}()

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Mock test: ", firstLog.Message)
	suite.Equal("got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestPanicErrorHandler_PanicStringError_ShouldMatchExpected() {
	func() {
		defer PanicErrorHandler(suite.logger, "Mock test: ")
		panic("got error")
	}()

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Mock test: ", firstLog.Message)
}

func (suite *HandlerSuite) TestPanicErrorHandler_PanicNoErrorNoStringTypeError_ShouldMatchExpected() {
	func() {
		defer PanicErrorHandler(suite.logger, "Mock test: ")
		panic(struct {
			Title string
		}{
			Title: "got error",
		})
	}()

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Mock test: ", firstLog.Message)
	suite.Equal("{Title:got error}", fmt.Sprintf("%+v", firstLog.Context[0].Interface))
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrValidate_ShouldStatusCodeGetStatusBadRequest() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrValidate, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusBadRequest, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[VALIDATE FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrVariable_ShouldStatusCodeGetStatusBadRequest() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrVariable, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusBadRequest, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[VARIABLE TYPE FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrInvalidArguments_ShouldStatusCodeGetStatusBadRequest() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrInvalidArguments, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusBadRequest, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[INVALID ARGUMENTS]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrAuthenticate_ShouldStatusCodeGetStatusUnauthorized() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrAuthenticate, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusUnauthorized, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[AUTHENTICATE FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrPermissionDeny_ShouldStatusCodeGetStatusForbidden() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrPermissionDeny, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusForbidden, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[PERMISSION DENY]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrJWTExecute_ShouldStatusCodeGetStatusForbidden() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrJWTExecute, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusForbidden, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[JWT EXECUTE FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrDatabaseRowNotFound_ShouldStatusCodeGetStatusNotFound() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrDatabaseRowNotFound, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusNotFound, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[DATABASE ROW NOT FOUND]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrExecute_ShouldStatusCodeGetStatusUnprocessableEntity() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrExecute, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusUnprocessableEntity, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[EXECUTE FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrJSONMarshal_ShouldStatusCodeGetStatusInternalServerError() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrJSONMarshal, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusInternalServerError, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[JSON MARSHAL FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrJSONUnmarshal_ShouldStatusCodeGetStatusInternalServerError() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrJSONUnmarshal, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusInternalServerError, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[JSON UNMARSHAL FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrDatabaseConnection_ShouldStatusCodeGetStatusServiceUnavailable() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrDatabaseConnection, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusServiceUnavailable, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[DATABASE CONNECTION FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassErrDatabaseDisconnect_ShouldStatusCodeGetStatusServiceUnavailable() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrDatabaseDisconnect, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusServiceUnavailable, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[DATABASE DISCONNECT FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func (suite *HandlerSuite) TestGinPanicErrorHandler_PassStringError_ShouldStatusCodeGetStatusServiceUnavailable() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler(suite.logger, "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		PanicIfErr(errors.New("got error"), ErrDatabaseDisconnect, errors.New("test error subject"))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusServiceUnavailable, result.StatusCode)

	suite.Equal(1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("[DATABASE DISCONNECT FAILED]: test error subject: got error", firstLog.Context[0].Interface.(error).Error())
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
