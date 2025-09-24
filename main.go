// @title User API
// @version 1.0
// @description This is a user management API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
	_ "github.com/lamkn06/user-app-golang.git/docs" // This is required for swagger
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	runtimeConfig runtime.ServerConfig
	dbConfig      runtime.DatabaseConfig
	jwtConfig     runtime.JWTConfig
)

type Server struct {
	config  runtime.ServerConfig
	routers []route.Router
	logger  *zap.SugaredLogger
}

func (s *Server) start() {
	server := echo.New()

	server.HTTPErrorHandler = middleware.ErrorHandler

	// Add Swagger route
	server.GET("/swagger/*", echoSwagger.WrapHandler)

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
	runtime.LoadConfigs([]any{&runtimeConfig, &dbConfig, &jwtConfig})

	logging.Init()
	logger := logging.NewSugaredLogger("server")
	logger.Infow("starting app")

	// Create startup context
	ctx := logging.AddLoggerToContext(context.Background(), logger)

	db, _ := repository.NewBunDB(ctx, dbConfig.PrimaryConnectionString())

	routers, err := route.Routers(ctx, runtimeConfig, db, jwtConfig)
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
