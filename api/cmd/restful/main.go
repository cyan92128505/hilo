package main

import (
	"hilo-api/pkg/errorCatcher"
	"hilo-api/pkg/shutdown"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer errorCatcher.PanicErrorHandler(logger, "Restful server interrupt => \n")

	_, cleanup, err := RestfulRunner()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	quit := make(chan os.Signal)
	defer close(quit)
	shutdown.NewShutdown(
		shutdown.WithQuit(quit),
	).Shutdown()
}
