package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lamkn06/user-app-golang.git/internal/repository"
	"github.com/lamkn06/user-app-golang.git/internal/route"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/pkg/logging"
	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
)

var (
	runtimeConfig runtime.ServerConfig
	dbConfig      runtime.DatabaseConfig
)

type Server struct {
	config  runtime.ServerConfig
	routers []route.Router
	logger  zerolog.Logger
}

func (s *Server) start() {
	server := echo.New()

	for _, r := range s.routers {
		r.Configure(server)
	}

	channel := make(chan error)
	go func() {
		channel <- server.Start(":" + s.config.Port)
	}()

	s.logger.Info().Msgf("Server started on port %s", s.config.Port)

	select {
	case sig := <-shutdownSignals():
		s.logger.Info().Msgf("Shutting down server... Received signal: %v", sig)
	case err := <-channel:
		s.logger.Error().Msgf("Failed to start server: %v", err)
	}

	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		s.logger.Error().Msgf("Failed to shutdown server: %v", err)
	}

	s.logger.Info().Msg("Server shutdown complete")
}

func main() {
	ctx := context.Background()
	logger := logging.NewLogger()

	ctx = logging.AddLoggerToContext(ctx, logger)
	runtime.LoadConfigs([]any{&runtimeConfig, &dbConfig})

	db, _ := repository.NewBunDB(ctx, dbConfig.PrimaryConnectionString())

	routers, err := route.Routers(ctx, runtimeConfig, db)
	if err != nil {
		logger.Error().Msgf("Failed to get routers: %v", err)
	}

	s := Server{routers: routers, config: runtimeConfig, logger: logger}
	s.start()
}

func shutdownSignals() (signals <-chan os.Signal) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGABRT, syscall.SIGTERM)
	return channel
}
