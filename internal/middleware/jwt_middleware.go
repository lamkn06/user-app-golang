package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/internal/service"
	"github.com/lamkn06/user-app-golang.git/pkg/exception"
)

func JWTMiddleware(jwtService service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				appErr := &exception.ApplicationError{
					Code:    exception.ErrorCodeUnauthorized,
					Message: "Authorization header required",
					Details: []exception.ErrorDetail{},
				}
				return c.JSON(appErr.HTTPStatus(), appErr)
			}

			// Extract token from "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				appErr := &exception.ApplicationError{
					Code:    exception.ErrorCodeUnauthorized,
					Message: "Invalid authorization header format",
					Details: []exception.ErrorDetail{},
				}
				return c.JSON(appErr.HTTPStatus(), appErr)
			}

			tokenString := tokenParts[1]
			token, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				appErr := &exception.ApplicationError{
					Code:    exception.ErrorCodeUnauthorized,
					Message: "Invalid token",
					Details: []exception.ErrorDetail{},
				}
				return c.JSON(appErr.HTTPStatus(), appErr)
			}

			// Extract user ID and add to context
			userID, err := jwtService.ExtractUserID(token)
			if err != nil {
				appErr := &exception.ApplicationError{
					Code:    exception.ErrorCodeUnauthorized,
					Message: "Invalid token claims",
					Details: []exception.ErrorDetail{},
				}
				return c.JSON(appErr.HTTPStatus(), appErr)
			}

			// Add user ID to context
			c.Set("userID", userID)
			return next(c)
		}
	}
}
