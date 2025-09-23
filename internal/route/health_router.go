package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/config"
	"github.com/lamkn06/user-app-golang.git/internal/response"
)

type HealthRouter struct {
}

func (r *HealthRouter) Configure(e *echo.Echo) {
	e.GET(config.GetVersionedAPIPath("/health"), r.HealthCheck)
}

func NewHealthRouter() *HealthRouter {
	return &HealthRouter{}
}

// HealthCheck godoc
//
// @Summary     Return API health status
// @Description Check if API status is ok
// @Produce     json
// @Success     200 {object} response.HealthResponse
// @Failure     500 {object} response.ErrorResponse
// @Router      /api/v1/health [get]
func (r *HealthRouter) HealthCheck(c echo.Context) (err error) {
	healthResp := response.HealthResponse{
		Status:  "OK",
		Version: config.APIVersion,
		Message: "API is healthy",
	}
	return c.JSON(http.StatusOK, healthResp)
}
