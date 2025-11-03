package jwt

import (
	"hilo-api/pkg/config"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
)

type testClaims struct {
	Email string `json:"email,omitempty"`
	*jwt.Claims
}

type JWTSuite struct {
	suite.Suite
	key    string
	option config.JWT
	jwt    IJWT
}

func (suite *JWTSuite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{
		PrivateKey: `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOChaSphj1MdLSxvU56h9vwmmpqdsQQF2alVwLKTj7dMoAoGCCqGSM49
AwEHoUQDQgAE7gMib5EUeW1An5VkkY4aU3xy+altlU3U0zn3FCO9Ffe/wwNUcUzp
XC9HWu76KhJnPpHczvZZv7Rro+kmqvN5tw==
-----END EC PRIVATE KEY-----`,
	}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
	j, err := NewES256JWT(result.PrivateKey)
	suite.NoError(err)
	suite.jwt = j
}

func (suite *JWTSuite) TestParseUnverified() {
	testEmail := "testMock@mock.com"
	standClaims := NewClaimsBuilder().
		WithSubject("testTopic").
		WithIssuer("tester").
		WithID("test001").
		WithAudience([]string{"testerClient"}).
		ExpiresAfter(5 * time.Second).Build()
	testClaimsInput := &testClaims{
		Claims: standClaims,
		Email:  testEmail,
	}
	result, err := suite.jwt.GenerateToken(testClaimsInput)
	suite.NoError(err)
	testClaimsResult := &testClaims{
		Claims: NewClaimsBuilder().Build(),
	}
	suite.NoError(ParseUnverified(result, testClaimsResult))
	suite.Equal(testEmail, testClaimsResult.Email)
}

func (suite *JWTSuite) TestParseUnverifiedError() {
	testClaimsResult := &testClaims{
		Claims: NewClaimsBuilder().Build(),
	}
	suite.Error(ParseUnverified("testToken", testClaimsResult))
}

func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTSuite))
}
