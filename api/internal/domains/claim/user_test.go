package claim

import (
	"encoding/json"
	"hilo-api/pkg/jwt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
}

func (suite *UserSuite) TestNewUser() {
	uid, err := uuid.NewRandom()
	suite.NoError(err)
	tk := NewUser(
		jwt.NewClaimsBuilder().WithSubject("testSubject").WithIssuer("testIssuer").ExpiresAfter(100*time.Second).Build(),
		WithUserID(uid.String()),
		WithPermissions("testPermission"),
	)

	suite.Equal("*claim.User", reflect.TypeOf(tk).String())
	suite.Equal([]string{"testPermission"}, tk.Permissions)
	result, err := json.Marshal(tk)
	suite.NoError(err)
	suite.T().Log(string(result))
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
