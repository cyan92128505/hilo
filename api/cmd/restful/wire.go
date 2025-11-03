// cmd/restful/wire.go
//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"fmt"
	"hilo-api/internal/domains/definition"
	restfulRouter "hilo-api/internal/presentation/restful"
	"hilo-api/pkg/config"
	"hilo-api/pkg/database/postgres"
	"hilo-api/pkg/jwt"
	"hilo-api/pkg/logger"
	"hilo-api/pkg/restful"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func ctx() context.Context {
	return context.Background()
}

var ctxSet = wire.NewSet(ctx)
var LoggerSet = wire.NewSet(logger.NewZap)

type Empty struct{}

func RunRestfulServer(logger *zap.Logger, coreOptions config.Set, route *gin.Engine, commonHandler restful.CommonHandler, handlers restfulRouter.HandlerSet) (Empty, func(), error) {
	restfulRouter.AddRoutes(route, commonHandler, handlers)
	if !coreOptions.Core.IsReleaseMode {
		pprof.Register(route)
	}

	h2s := &http2.Server{}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", coreOptions.Server.Port),
		Handler: h2c.NewHandler(route, h2s),
	}

	go func(s *http.Server) {
		logger.Info("start restful server",
			zap.String("system", coreOptions.Core.SystemName),
			zap.String("port", coreOptions.Server.Port),
		)
		if err := s.ListenAndServe(); err != nil {
			logger.Warn("restful server error or closed",
				zap.String("system", coreOptions.Core.SystemName),
				zap.Error(err),
			)
		}
	}(httpServer)
	return Empty{}, func() {
		ctx, cancel := context.WithTimeout(context.Background(), coreOptions.Server.ServerTimeout)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Warn("restful server Failed to Shutdown",
				zap.String("system", coreOptions.Core.SystemName),
				zap.Error(err),
			)
		}
	}, nil
}

func RestfulRunner() (Empty, func(), error) {
	panic(wire.Build(wire.NewSet(
		ctxSet,
		wire.NewSet(
			config.NewSet,
			config.NewCore,
			config.NewJWT,
			config.NewPostgres,
			config.NewServer,
		),
		LoggerSet,
		postgres.NewPostgresDB,
		// wire.NewSet(
		// 	wire.Struct(new(repository.Set), "*")),
		wire.NewSet(jwt.NewES256JWTFromOptions, wire.Bind(new(definition.ES256JWT), new(*jwt.ES256JWT))),
		// wire.NewSet(
		// 	wire.Struct(new(usecase.Set), "*")),
		wire.NewSet(restfulRouter.NewAPIGuardValidator, wire.Bind(new(restful.GuarderValidator), new(*restfulRouter.APIGuardValidator))),
		wire.NewSet(restful.NewJWTGuarder),
		wire.NewSet(restful.NewGin),
		wire.Value(restful.CommonHandler{
			Error404:   restful.Error404Set,
			QuickReply: restful.QuickReplySet,
			PromHTTP:   restful.NewPromHTTPSet,
		}),
		wire.NewSet(
			wire.Struct(new(restfulRouter.HandlerSet), "*")),
		RunRestfulServer,
	)))
}
