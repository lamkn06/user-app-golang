package route

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Router interface {
	// Configure configures the router
	Configure(e *echo.Echo)
}

func Routers(ctx context.Context) (routers []Router, err error) {

	return []Router{
		NewHealthRouter(),
	}, err
}
