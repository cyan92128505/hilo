package errorCatcher

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *ApiError
		if errors.As(err, &appErr) {
			// Handle ApiError
			logger.Error("request failed",
				zap.String("code", appErr.Code),
				zap.String("message", appErr.Message),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(appErr.Cause),
			)

			c.JSON(appErr.StatusCode, gin.H{
				"error":   appErr.Code,
				"message": appErr.Message,
			})
			return
		}

		// Handle unknown errors
		logger.Error("unexpected error",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Error(err),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "INTERNAL_ERROR",
			"message": "Internal server error",
		})
	}
}
