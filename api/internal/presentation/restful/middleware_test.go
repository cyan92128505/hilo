package restful

import (
	"hilo-api/internal/domains/claim"
	"hilo-api/pkg/config"
	"hilo-api/pkg/jwt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type APIGuardValidatorSuite struct {
	suite.Suite
	c                   *gin.Context
	e                   *gin.Engine
	jwtOp               config.JWT
	jwt                 jwt.IJWT
	token               string
	expiredToken        string
	noneToken           string
	nonePermissionToken string
}

func (suite *APIGuardValidatorSuite) SetupSuite() {
	suite.c, suite.e = gin.CreateTestContext(httptest.NewRecorder())
	suite.c.Request = &http.Request{
		RequestURI: "/ping",
	}

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

	token, err := suite.jwt.GenerateToken(claim.NewUser(
		jwt.NewClaimsBuilder().ExpiresAfter(500*time.Second).Build(),
		claim.WithPermissions("/ping"),
	))
	suite.NoError(err)
	suite.token = token

	expiredToken, err := suite.jwt.GenerateToken(claim.NewUser(
		jwt.NewClaimsBuilder().ExpiresAfter(-50*time.Second).Build(),
		claim.WithPermissions("/ping"),
	))
	suite.NoError(err)
	suite.expiredToken = expiredToken

	noneToken, err := suite.jwt.GenerateToken(claim.NewUser(
		jwt.NewClaimsBuilder().Build(),
		claim.WithPermissions("/ping"),
	))
	suite.NoError(err)
	suite.noneToken = noneToken

	nonePermissionToken, err := suite.jwt.GenerateToken(claim.NewUser(
		jwt.NewClaimsBuilder().Build(),
	))
	suite.NoError(err)
	suite.nonePermissionToken = nonePermissionToken
}

func (suite *APIGuardValidatorSuite) TestNewAPIGuardValidator() {
	suite.Equal("*restful.APIGuardValidator", reflect.TypeOf(NewAPIGuardValidator(suite.jwt)).String())
}

func (suite *APIGuardValidatorSuite) TestVerify() {
	suite.NoError(NewAPIGuardValidator(suite.jwt).Verify(suite.c, suite.token))
}

func (suite *APIGuardValidatorSuite) TestVerifyExpired() {
	suite.Error(NewAPIGuardValidator(suite.jwt).Verify(suite.c, suite.expiredToken))
}

func (suite *APIGuardValidatorSuite) TestVerifyNoExpired() {
	suite.NoError(NewAPIGuardValidator(suite.jwt).Verify(suite.c, suite.noneToken))
}

func (suite *APIGuardValidatorSuite) TestVerifyNoPermission() {
	suite.Error(NewAPIGuardValidator(suite.jwt).Verify(suite.c, suite.nonePermissionToken))
}

func TestAPIGuardValidatorSuite(t *testing.T) {
	suite.Run(t, new(APIGuardValidatorSuite))
}
