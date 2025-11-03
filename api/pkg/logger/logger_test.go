package logger

import (
	"errors"
	"hilo-api/pkg/config"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type LoggerTestSuite struct {
	suite.Suite
}

func (suite *LoggerTestSuite) TestNewZap() {
	logger, err := NewZap(config.Core{
		LogLevel:      "ERROR",
		IsReleaseMode: false,
	})
	suite.NoError(err)
	suite.NotNil(logger)
	suite.Equal(logger, zap.L())
}

func (suite *LoggerTestSuite) TestPanicLogger() {
	logger, err := NewZap(config.Core{
		LogLevel:      "WARN",
		IsReleaseMode: true,
	})
	suite.NoError(err)
	func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Info("INFO message", zap.Error(errors.New("test info error")))
				logger.Warn("WARN message", zap.Error(errors.New("test warn error")))
				logger.Error("ERROR message", zap.Error(errors.New("test error error")))
			}
		}()
		panic("test panic")
	}()
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}
