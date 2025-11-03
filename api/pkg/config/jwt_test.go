package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JWTSuite struct {
	suite.Suite
	PrivateKeyPath string
	PrivateKey     string
}

func (suite *JWTSuite) SetupSuite() {
	os.Clearenv()
	suite.PrivateKeyPath = "testEcdsaPrivateKeyPath"
	suite.PrivateKey = "testPrivateKey"

	suite.NoError(os.Setenv("PRIVATE_KEY_PATH", suite.PrivateKeyPath))
	suite.NoError(os.Setenv("PRIVATE_KEY", suite.PrivateKey))
}

func (suite *JWTSuite) TestDefaultOption() {
	jwt := &JWT{}
	suite.NoError(LoadFromEnv(jwt))
	suite.Equal(suite.PrivateKeyPath, jwt.PrivateKeyPath)
	suite.Equal(suite.PrivateKey, jwt.PrivateKey)
}

func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTSuite))
}
