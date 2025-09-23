package route

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
)

type Router interface {
	// Configure configures the router
	Configure(e *echo.Echo)
}

func Routers(ctx context.Context, config runtime.ServerConfig) (routers []Router, err error) {
	return []Router{
		NewHealthRouter(config),
	}, err
}
