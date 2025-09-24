package route

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/middleware"
	"github.com/lamkn06/user-app-golang.git/internal/runtime"
	"github.com/lamkn06/user-app-golang.git/internal/service"
	"github.com/lamkn06/user-app-golang.git/pkg/api/request"
	"github.com/lamkn06/user-app-golang.git/pkg/exception"
	"github.com/lamkn06/user-app-golang.git/pkg/logging"
)

type AuthRouter struct {
	config      runtime.ServerConfig
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthRouter(config runtime.ServerConfig, authService service.AuthService) *AuthRouter {
	return &AuthRouter{
		config:      config,
		authService: authService,
		validator:   validator.New(),
	}
}

func (r *AuthRouter) Configure(e *echo.Echo) {
	e.POST("/api/"+r.config.APIVersion+"/auth/signup", r.SignUp)
	e.POST("/api/"+r.config.APIVersion+"/auth/signin", r.SignIn)
	e.POST("/api/"+r.config.APIVersion+"/auth/signout", r.SignOut)
}

// SignUp godoc
// @Summary Sign up a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body request.SignUpRequest true "User registration information"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /auth/signup [post]
func (r *AuthRouter) SignUp(c echo.Context) error {
	logger := logging.LoggerFromContext(c.Request().Context())
	req := request.SignUpRequest{}

	if err := c.Bind(&req); err != nil {
		appErr := exception.ToApplicationError(err, exception.ErrorCodeBadRequest)
		logger.Errorw("Failed to bind signup request", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	if err := r.validator.Struct(req); err != nil {
		appErr := middleware.ParseValidationError(err)
		logger.Errorw("Failed to validate signup request", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	authResp, err := r.authService.SignUp(req)
	if err != nil {
		logger.Errorw("Failed to sign up user", "error", err)
		appErr := exception.ToApplicationError(err, exception.ErrorCodeInternalServerError)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	return c.JSON(http.StatusOK, authResp)
}

// SignIn godoc
// @Summary Sign in a user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body request.SignInRequest true "User credentials"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} exception.ApplicationError
// @Failure 401 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /auth/signin [post]
func (r *AuthRouter) SignIn(c echo.Context) error {
	logger := logging.LoggerFromContext(c.Request().Context())
	req := request.SignInRequest{}

	if err := c.Bind(&req); err != nil {
		appErr := exception.ToApplicationError(err, exception.ErrorCodeBadRequest)
		logger.Errorw("Failed to bind signin request", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	if err := r.validator.Struct(req); err != nil {
		appErr := middleware.ParseValidationError(err)
		logger.Errorw("Failed to validate signin request", "error", err)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	authResp, err := r.authService.SignIn(req)
	if err != nil {
		logger.Errorw("Failed to sign in user", "error", err)
		appErr := &exception.ApplicationError{
			Code:    exception.ErrorCodeUnauthorized,
			Message: "Invalid credentials",
			Details: []exception.ErrorDetail{},
		}
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	return c.JSON(http.StatusOK, authResp)
}

// SignOut godoc
// @Summary Sign out a user
// @Description Sign out user and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SignOutResponse
// @Failure 401 {object} exception.ApplicationError
// @Failure 500 {object} exception.ApplicationError
// @Router /auth/signout [post]
func (r *AuthRouter) SignOut(c echo.Context) error {
	logger := logging.LoggerFromContext(c.Request().Context())

	// Get token from Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		appErr := &exception.ApplicationError{
			Code:    exception.ErrorCodeUnauthorized,
			Message: "Authorization header required",
			Details: []exception.ErrorDetail{},
		}
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	signOutResp, err := r.authService.SignOut(authHeader)
	if err != nil {
		logger.Errorw("Failed to sign out user", "error", err)
		appErr := exception.ToApplicationError(err, exception.ErrorCodeInternalServerError)
		return c.JSON(appErr.HTTPStatus(), appErr)
	}

	return c.JSON(http.StatusOK, signOutResp)
}
