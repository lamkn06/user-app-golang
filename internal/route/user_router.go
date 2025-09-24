package route

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/middleware"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/internal/service"
	"github.com/lamkn06/user-app-golang.git/pkg/api/request"
	"github.com/lamkn06/user-app-golang.git/pkg/exception"
	"github.com/lamkn06/user-app-golang.git/pkg/logging"
)

type UserRouter struct {
	config      runtime.ServerConfig
	userService service.UserService
	validator   *validator.Validate
}

func NewUserRouter(config runtime.ServerConfig, userService service.UserService) *UserRouter {
	return &UserRouter{config: config, userService: userService, validator: validator.New()}
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
	logger := logging.LoggerFromContext(c.Request().Context())
	user := request.NewUserRequest{}

	if err := c.Bind(&user); err != nil {
		appErr := exception.ToApplicationError(err, exception.ErrorCodeBadRequest)
		logger.Errorw("Failed to bind user", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	if err := r.validator.Struct(user); err != nil {
		appErr := middleware.ParseValidationError(err)
		logger.Errorw("Failed to validate user", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	newUser, err := r.userService.NewUser(user)
	if err != nil {
		logger.Errorw("Failed to create user", "error", err)
		appErr := exception.ToApplicationError(err, exception.ErrorCodeInternalServerError)
		return c.JSON(appErr.HTTPStatus(), appErr)
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
