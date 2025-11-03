package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	suite.Suite
	PostgresURL       string
	PostgresTxTimeout time.Duration
}

func (suite *PostgresSuite) SetupSuite() {
	os.Clearenv()
	suite.PostgresURL = "testPostgresURL"
	suite.PostgresTxTimeout = 20 * time.Second
	suite.NoError(os.Setenv("POSTGRES_URL", suite.PostgresURL))
	suite.NoError(os.Setenv("POSTGRES_TX_TIMEOUT", fmt.Sprint(suite.PostgresTxTimeout)))
}

func (suite *PostgresSuite) TestDefaultOption() {
	postgres := &Postgres{}
	suite.NoError(LoadFromEnv(postgres))
	suite.Equal(suite.PostgresURL, postgres.PostgresURL)
	suite.Equal(suite.PostgresTxTimeout, postgres.PostgresTxTimeout)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
