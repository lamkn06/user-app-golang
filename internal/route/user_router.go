package route

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/internal/service"
	"github.com/lamkn06/user-app-golang.git/pkg/api/request"
)

type UserRouter struct {
	config      runtime.ServerConfig
	userService service.UserService
}

func NewUserRouter(config runtime.ServerConfig, userService service.UserService) *UserRouter {
	return &UserRouter{config: config, userService: userService}
}

func (r *UserRouter) Configure(e *echo.Echo) {
	e.GET("/api/"+r.config.APIVersion+"/users", r.GetUsers)
	e.POST("/api/"+r.config.APIVersion+"/users", r.CreateUser)
	e.GET("/api/"+r.config.APIVersion+"/users/:id", r.GetUserById)
}

func (r *UserRouter) GetUsers(c echo.Context) error {
	users, err := r.userService.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, users)
}

func (r *UserRouter) CreateUser(c echo.Context) error {

	user := request.NewUserRequest{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	newUser, err := r.userService.NewUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, newUser)
}

func (r *UserRouter) GetUserById(c echo.Context) error {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := r.userService.GetUserById(idUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}
