package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSetSuite struct {
	suite.Suite
}

func (suite *ConfigSetSuite) TestNewConfigSet() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("config.Set", reflect.TypeOf(result).String())
}

func (suite *ConfigSetSuite) TestNewCore() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Core", reflect.TypeOf(NewCore(result)).Name())
}

func (suite *ConfigSetSuite) TestNewJWT() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("JWT", reflect.TypeOf(NewJWT(result)).Name())
}

func (suite *ConfigSetSuite) TestNewPostgres() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Postgres", reflect.TypeOf(NewPostgres(result)).Name())
}

func (suite *ConfigSetSuite) TestNewServer() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Server", reflect.TypeOf(NewServer(result)).Name())
}

func TestConfigSetSuite(t *testing.T) {
	suite.Run(t, new(ConfigSetSuite))
}
