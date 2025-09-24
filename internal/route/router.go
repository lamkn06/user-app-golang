package route

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/repository"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/internal/service"
	"github.com/uptrace/bun"
)

type Router interface {
	Configure(e *echo.Echo)
}

func Routers(ctx context.Context, config runtime.ServerConfig, db *bun.DB, jwtConfig runtime.JWTConfig) (routers []Router, err error) {
	userRepository := repository.NewUserRepository(db, ctx)
	userService := service.NewUserService(userRepository)

	jwtService := service.NewJWTService(jwtConfig)
	authService := service.NewAuthService(userRepository, jwtService)

	return []Router{
		NewHealthRouter(config),
		NewUserRouter(config, userService, jwtService),
		NewAuthRouter(config, authService),
	}, err
}
