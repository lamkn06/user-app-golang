package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/pkg/api/response"
	"github.com/lamkn06/user-app-golang.git/pkg/exception"
)

func unwrapRecursive(err error) error {
	var originalErr = err

	for originalErr != nil {
		var internalErr = errors.Unwrap(originalErr)

		if internalErr == nil {
			break
		}

		originalErr = internalErr
	}

	return originalErr
}

func ErrorHandler(err error, c echo.Context) {
	var (
		statusCode    int
		errorResponse interface{}
	)

	switch e := err.(type) {
	case *exception.ApplicationError:
		statusCode = http.StatusBadRequest
		errorResponse = toErrorResponse(*e)
	default:
		var he *echo.HTTPError
		ok := errors.As(err, &he)
		if !ok {
			he = &echo.HTTPError{
				Code: http.StatusInternalServerError,
				Message: toErrorResponse(exception.ApplicationError{
					Message: unwrapRecursive(err).Error(),
					Code:    exception.ErrorCodeInternalServerError,
				}),
			}
		}

		statusCode = he.Code
		errorResponse = he.Message
	}

	_err := c.JSON(statusCode, errorResponse)
	if _err != nil {
		c.Echo().Logger.Error(_err)
	}
}

func toErrorResponse(error exception.ApplicationError) response.ErrorResponse {
	return response.ErrorResponse{
		Message: error.Message,
		Code:    error.Code,
		Details: error.Details,
	}
}
