package restful

import (
	"errors"
	"hilo-api/internal/domain/claim"
	"hilo-api/internal/domain/definition"
	authDefinition "hilo-api/pkg/definition"
	"hilo-api/pkg/errorCatcher"
	jwtTool "hilo-api/pkg/jwt"
	"hilo-api/pkg/restful"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewAPIGuardValidator method
func NewAPIGuardValidator(jwt definition.ES256JWT) *APIGuardValidator {
	return &APIGuardValidator{
		jwt: jwt,
	}
}

// APIGuardValidator method
type APIGuardValidator struct {
	jwt definition.ES256JWT
}

// Verify method
func (b *APIGuardValidator) Verify(c *gin.Context, token string) error {
	userClaim := claim.NewUser(jwtTool.NewClaimsBuilder().Build())
	if err := b.jwt.VerifyToken(token, userClaim); err != nil {
		return errorCatcher.ConcatError(
			errorCatcher.ErrAuthenticate,
			restful.ErrValidatorVerify,
			err,
		)
	}
	// set user id to gin context
	c.Set(GinContextUserIDKey, userClaim.UserID)
	for _, permission := range userClaim.Permissions {
		if strings.HasPrefix(c.FullPath(), permission) || strings.HasPrefix(c.Request.RequestURI, permission) {
			c.Set(authDefinition.AuthorizationKey, userClaim)
			return nil
		}
	}
	return errorCatcher.ConcatError(
		errorCatcher.ErrPermissionDeny,
		restful.ErrValidatorVerify,
		errors.New("no permission allowed to resource"),
	)
}
