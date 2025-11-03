package restful

import (
	"errors"
	"hilo-api/pkg/definition"
	"hilo-api/pkg/errorCatcher"
	jwtTool "hilo-api/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrValidatorVerify = errors.New("[Validator Verify Failed]")
)

func NewBasicGuardValidator(jwt jwtTool.IJWT) *BasicGuardValidator {
	return &BasicGuardValidator{
		jwt: jwt,
	}
}

type BasicGuardValidator struct {
	jwt jwtTool.IJWT
}

func (b *BasicGuardValidator) Verify(c *gin.Context, token string) error {
	commonClaims := jwtTool.NewCommon(jwtTool.NewClaimsBuilder().Build())
	if err := b.jwt.VerifyToken(token, commonClaims); err != nil {
		return errorCatcher.ConcatError(
			errorCatcher.ErrAuthenticate,
			ErrValidatorVerify,
			err,
		)
	}

	for _, permission := range commonClaims.Permissions {
		if strings.HasPrefix(c.FullPath(), permission) || strings.HasPrefix(c.Request.RequestURI, permission) {
			c.Set(definition.AuthTokenKey, commonClaims)
			return nil
		}
	}
	return errorCatcher.ConcatError(
		errorCatcher.ErrPermissionDeny,
		ErrValidatorVerify,
		errors.New("no permission allowed to resource"),
	)
}
