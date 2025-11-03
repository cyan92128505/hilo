package restful

import (
	"hilo-api/pkg/config"
	"hilo-api/pkg/errorCatcher"
	"hilo-api/pkg/jwt"
	"hilo-api/pkg/logger"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type testGuarderValidator struct {
	mock.Mock
	GuarderValidator
}

func (t *testGuarderValidator) Verify(c *gin.Context, token string) error {
	args := t.Called(c, token)
	return args.Error(0)
}

type MiddlewareSuite struct {
	suite.Suite
	logger *zap.Logger
	jwtOp  config.JWT
	jwt    jwt.IJWT
}

func (suite *MiddlewareSuite) SetupSuite() {
	zapLogger, err := logger.NewZap(config.Core{
		SystemName: "testSystemName",
	})
	suite.NoError(err)
	suite.logger = zapLogger
	suite.jwtOp = config.JWT{
		PrivateKey: `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOChaSphj1MdLSxvU56h9vwmmpqdsQQF2alVwLKTj7dMoAoGCCqGSM49
AwEHoUQDQgAE7gMib5EUeW1An5VkkY4aU3xy+altlU3U0zn3FCO9Ffe/wwNUcUzp
XC9HWu76KhJnPpHczvZZv7Rro+kmqvN5tw==
-----END EC PRIVATE KEY-----
`}
	es256, err := jwt.NewES256JWT(suite.jwtOp.PrivateKey)
	suite.NoError(err)
	suite.jwt = es256

}

func (suite *MiddlewareSuite) TestNewJWTGuard() {
	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, mock.Anything).Return(nil)
	result := NewJWTGuarder(testGuarderValidator).JWTGuarder()
	suite.Equal("gin.HandlerFunc", reflect.TypeOf(result).String())
}

func (suite *MiddlewareSuite) TestJWTGuarderRun() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(
		jwt.NewClaimsBuilder().ExpiresAfter(500*time.Second).Build(),
		jwt.WithPermissions("/ping"),
	))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(nil)

	r := gin.Default()
	r.GET("/ping", NewJWTGuarder(testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunWhiteList() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(
		jwt.NewClaimsBuilder().ExpiresAfter(500*time.Second).Build(),
		jwt.WithPermissions("/ping"),
	))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(nil)

	r := gin.Default()
	r.GET("/ping", NewJWTGuarder(testGuarderValidator).JWTGuarder("/ping"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunNotAuthorization() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(jwt.NewClaimsBuilder().Build()))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}

	r := gin.Default()
	r.Use(gin.Logger(), errorCatcher.GinPanicErrorHandler(suite.logger, "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Auth", "Basic "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunFormatError() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(jwt.NewClaimsBuilder().Build()))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(nil)

	r := gin.Default()
	r.Use(gin.Logger(), errorCatcher.GinPanicErrorHandler(suite.logger, "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Basic "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunExpired() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(jwt.NewClaimsBuilder().ExpiresAfter(1 * time.Second).Build()))
	suite.NoError(err)

	time.Sleep(1 * time.Second)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(errorCatcher.ErrPermissionDeny)

	r := gin.Default()
	r.Use(gin.Logger(), errorCatcher.GinPanicErrorHandler(suite.logger, "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunPermissionNotAllow() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(
		jwt.NewClaimsBuilder().ExpiresAfter(500 * time.Second).Build(),
	))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(errorCatcher.ErrPermissionDeny)

	r := gin.Default()
	r.Use(gin.Logger(), errorCatcher.GinPanicErrorHandler(suite.logger, "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
