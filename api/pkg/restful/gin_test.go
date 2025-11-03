package restful

import (
	"hilo-api/pkg/config"
	"hilo-api/pkg/logger"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type GinSuite struct {
	suite.Suite
	originCoreOption    config.Core
	originServerOption  config.Server
	anotherCoreOption   config.Core
	anotherServerOption config.Server
	logger              *zap.Logger
}

func (suite *GinSuite) SetupSuite() {
	suite.originCoreOption = config.Core{
		SystemName:    "Gin Server",
		IsReleaseMode: true,
		LogLevel:      "warn",
	}
	suite.originServerOption = config.Server{
		PrefixMessage:    "error gin server",
		AllowAllOrigins:  true,
		CustomizedRender: true,
	}

	suite.anotherCoreOption = config.Core{
		SystemName: "Gin Server",
	}

	suite.anotherServerOption = config.Server{
		ReleaseMode:     true,
		PrefixMessage:   "error gin server",
		AllowAllOrigins: false,
		AllowOrigins:    []string{"http://localhost"},
	}
	zapLogger, err := logger.NewZap(suite.originCoreOption)
	suite.NoError(err)
	suite.logger = zapLogger
}

func (suite *GinSuite) TestNewGin() {
	gin, err := NewGin(suite.logger, suite.originServerOption, &JWTGuarder{})
	suite.NoError(err)
	suite.Equal("*gin.Engine", reflect.TypeOf(gin).String())
}

func (suite *GinSuite) TestNewGinAllowOrigins() {
	gin, err := NewGin(suite.logger, suite.originServerOption, &JWTGuarder{})
	suite.NoError(err)
	suite.Equal("*gin.Engine", reflect.TypeOf(gin).String())
}

func (suite *GinSuite) TestNewGinAllowOriginsReleaseAndLimitOrigin() {
	gin, err := NewGin(suite.logger, suite.anotherServerOption, &JWTGuarder{})
	suite.NoError(err)
	suite.Equal("*gin.Engine", reflect.TypeOf(gin).String())
}

func TestGinSuite(t *testing.T) {
	suite.Run(t, new(GinSuite))
}
