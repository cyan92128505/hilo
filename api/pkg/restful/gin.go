package restful

import (
	"hilo-api/pkg/config"
	"hilo-api/pkg/errorCatcher"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewGin(
	logger *zap.Logger,
	cfgServer config.Server,
	guarder *JWTGuarder,
) (*gin.Engine, error) {
	if cfgServer.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	srv := gin.New()

	srv.MaxMultipartMemory = cfgServer.MaxMultipartMemoryMB << 20

	cf := cors.DefaultConfig()
	cf.AllowMethods = []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
	}
	cf.AllowHeaders = []string{
		"Origin",
		"UpgradePost",
		"Upgrade",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Connection",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Host",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
		"X-Requested-With",
		"X-Google-*",
		"X-AppEngine-*",
		"X-CloudScheduler",
		"X-CloudScheduler-JobName",
		"X-CloudScheduler-ScheduleTime",
		"Sec-WebSocket-Key",
		"Sec-WebSocket-Version",
		"Sec-WebSocket-Protocol",
	}

	cf.ExposeHeaders = []string{
		"Content-Length",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
	}

	if cfgServer.AllowAllOrigins {
		cf.AllowAllOrigins = true
	} else {
		cf.AllowOrigins = cfgServer.AllowOrigins
	}
	fns := []gin.HandlerFunc{
		cors.New(cf),
		gin.Logger(),
		errorCatcher.GinPanicErrorHandler(logger, cfgServer.PrefixMessage),
	}
	if cfgServer.JWTGuard {
		fns = append(fns, guarder.JWTGuarder(cfgServer.AllowedPaths...))
	}
	srv.Use(fns...)

	return srv, nil
}
