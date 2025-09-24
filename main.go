package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lamkn06/user-app-golang.git/internal/middleware"
	"github.com/lamkn06/user-app-golang.git/internal/repository"
	"github.com/lamkn06/user-app-golang.git/internal/route"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/pkg/logging"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

var (
	runtimeConfig runtime.ServerConfig
	dbConfig      runtime.DatabaseConfig
)

type Server struct {
	config  runtime.ServerConfig
	routers []route.Router
	logger  *zap.SugaredLogger
}

func (s *Server) start() {
	server := echo.New()

	server.HTTPErrorHandler = middleware.ErrorHandler

	for _, r := range s.routers {
		r.Configure(server)
	}

	channel := make(chan error)
	go func() {
		channel <- server.Start(":" + s.config.Port)
	}()

	s.logger.Infof("Server started on port %s", s.config.Port)

	select {
	case sig := <-shutdownSignals():
		s.logger.Infof("Shutting down server... Received signal: %v", sig)
	case err := <-channel:
		s.logger.Errorf("Failed to start server: %v", err)
	}

	ctx := context.Background()
	if err := server.Shutdown(ctx); err != nil {
		s.logger.Errorf("Failed to shutdown server: %v", err)
	}

	s.logger.Info("Server shutdown complete")
}

func main() {
	runtime.LoadConfigs([]any{&runtimeConfig, &dbConfig})

	logging.Init()
	logger := logging.NewSugaredLogger("server")
	logger.Infow("starting app")

	// Create startup context
	ctx := logging.AddLoggerToContext(context.Background(), logger)

	db, _ := repository.NewBunDB(ctx, dbConfig.PrimaryConnectionString())

	routers, err := route.Routers(ctx, runtimeConfig, db)
	if err != nil {
		logger.Errorw("Failed to get routers", "error", err)
	}

	s := Server{routers: routers, config: runtimeConfig, logger: logger}
	s.start()
}

func shutdownSignals() (signals <-chan os.Signal) {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGABRT, syscall.SIGTERM)
	return channel
}
