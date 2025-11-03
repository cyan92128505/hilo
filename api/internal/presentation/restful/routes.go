package restful

import (
	"hilo-api/pkg/restful"

	"github.com/gin-gonic/gin"
)

// HandlerSet struct
type HandlerSet struct {
}

// AddRoutes func
func AddRoutes(route *gin.Engine, commonHandler restful.CommonHandler, handlers HandlerSet) {
	route.GET("/ping", commonHandler.QuickReply)
	route.GET("/metrics", commonHandler.PromHTTP)

	route.NoRoute(commonHandler.Error404)
}
