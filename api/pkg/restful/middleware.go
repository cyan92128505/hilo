package restful

import (
	"errors"
	"hilo-api/pkg/definition"
	"hilo-api/pkg/errorCatcher"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrGuardJWTGuarder = errors.New("[Guard JWT Guarder Failed]")
)

type GuarderValidator interface {
	Verify(c *gin.Context, token string) error
}

func NewJWTGuarder(validator GuarderValidator) *JWTGuarder {
	return &JWTGuarder{
		validator: validator,
	}
}

type JWTGuarder struct {
	validator GuarderValidator
}

// JWTGuarder method
func (j *JWTGuarder) JWTGuarder(allowList ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip if matching allowList
		for _, term := range allowList {
			if strings.HasPrefix(c.FullPath(), term) || strings.HasPrefix(c.Request.RequestURI, term) {
				c.Next()
				return
			}
		}
		token := c.DefaultQuery(definition.QueryAuthKey, "")
		if token == "" {
			authorization := c.GetHeader(definition.AuthorizationKey)
			if len(authorization) == 0 {
				panic(
					errorCatcher.ConcatError(
						errorCatcher.ErrAuthenticate,
						ErrGuardJWTGuarder,
						errors.New("metadata key: [authorization] must required"),
					),
				)
			}
			if strings.HasPrefix(authorization, definition.AuthorizationType) {
				token = strings.TrimPrefix(authorization, definition.AuthorizationType)
			} else {
				panic(
					errorCatcher.ConcatError(
						errorCatcher.ErrAuthenticate,
						ErrGuardJWTGuarder,
						errors.New("JWT Authorization format error: must be Bearer"),
					),
				)
			}
		}

		if err := j.validator.Verify(c, token); err != nil {
			panic(err)
		}
	}
}
