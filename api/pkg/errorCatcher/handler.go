package errorCatcher

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PanicErrorHandler func
// only logger error message if panic
// no andy side effect or throw after catch error
func PanicErrorHandler(logger *zap.Logger, prefixMessage string) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			logger.Error(prefixMessage, zap.Error(e))
		default:
			logger.Error(prefixMessage, zap.Any("data", e))
		}
	}
}

func GinPanicErrorHandler(logger *zap.Logger, prefixMessage string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch e := err.(type) {
				case error:
					statusCode := http.StatusUnauthorized
					logger.Error(prefixMessage, zap.Error(e))
					switch {
					case errors.Is(e, ErrValidate),
						errors.Is(e, ErrVariable),
						errors.Is(e, ErrGinBindingAndValidate),
						errors.Is(e, ErrInvalidArguments):
						statusCode = http.StatusBadRequest
						break
					case errors.Is(e, ErrAuthenticate):
						statusCode = http.StatusUnauthorized
						break
					case errors.Is(e, ErrPermissionDeny),
						errors.Is(e, ErrJWTExecute):
						statusCode = http.StatusForbidden
						break
					case errors.Is(e, ErrDatabaseRowNotFound):
						statusCode = http.StatusNotFound
						break
					case errors.Is(e, ErrExecute),
						errors.Is(e, ErrDatabaseExecute),
						errors.Is(e, ErrDatabaseExecuteNotNullViolation),
						errors.Is(e, ErrDatabaseExecuteForeignKeyViolation),
						errors.Is(e, ErrDatabaseExecuteUniqueViolation),
						errors.Is(e, ErrDatabaseExecuteCheckViolation),
						errors.Is(e, ErrDatabaseExecuteMultipleColumnUpdateMustSubSelect):
						statusCode = http.StatusUnprocessableEntity
						break
					case errors.Is(e, ErrJSONMarshal),
						errors.Is(e, ErrJSONUnmarshal):
						statusCode = http.StatusInternalServerError
						break
					case errors.Is(e, ErrDatabaseConnection),
						errors.Is(e, ErrDatabaseDisconnect):
						statusCode = http.StatusServiceUnavailable
						break
					}
					c.AbortWithError(statusCode, fmt.Errorf("%s: %w", prefixMessage, err.(error)))
					break
				default:
					logger.Error(prefixMessage, zap.Any("data", e))
					c.AbortWithError(http.StatusInternalServerError, errors.New(prefixMessage))
				}
			}
		}()
		c.Next()
	}
}
