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
	jwtService  service.JWTService
	validator   *validator.Validate
}

func NewUserRouter(config runtime.ServerConfig, userService service.UserService, jwtService service.JWTService) *UserRouter {
	return &UserRouter{config: config, userService: userService, jwtService: jwtService, validator: validator.New()}
}

func (r *UserRouter) Configure(e *echo.Echo) {
	e.GET("/api/"+r.config.APIVersion+"/users", r.GetUsers)
	e.POST("/api/"+r.config.APIVersion+"/users", r.CreateUser)
	e.GET("/api/"+r.config.APIVersion+"/users/:id", r.GetUserById, middleware.JWTMiddleware(r.jwtService))
}

// GetUsers godoc
// @Summary Get all users with pagination
// @Description Get all users from the database with pagination support
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.ListResponse[response.NewUserResponse]
// @Failure 400 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /users [get]
func (r *UserRouter) GetUsers(c echo.Context) error {
	logger := logging.LoggerFromContext(c.Request().Context())

	// Parse pagination parameters with defaults
	listReq := request.NewListRequest()
	if err := c.Bind(&listReq); err != nil {
		appErr := exception.ToApplicationError(err, exception.ErrorCodeBadRequest)
		logger.Errorw("Failed to bind pagination parameters", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	// Validate pagination parameters
	if err := r.validator.Struct(listReq); err != nil {
		appErr := middleware.ParseValidationError(err)
		logger.Errorw("Failed to validate pagination parameters", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	// Get users with pagination
	users, err := r.userService.GetUsers(listReq)
	if err != nil {
		logger.Errorw("Failed to get users", "error", err)
		appErr := exception.ToApplicationError(err, exception.ErrorCodeInternalServerError)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	return c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.NewUserRequest true "User information"
// @Success 200 {object} response.NewUserResponse
// @Failure 400 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /users [post]
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
		return exception.ToApplicationError(err, exception.ErrorCodeInternalServerError)
	}

	return c.JSON(http.StatusOK, newUser)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.NewUserResponse
// @Failure 400 {object} exception.ApplicationError
// @Failure 401 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /users/{id} [get]
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
