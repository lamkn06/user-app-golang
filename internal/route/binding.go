package route

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lamkn06/user-app-golang.git/pkg/exception"
)

func BindAndValidate[T interface{}](c echo.Context, validate validator.Validate, requestData T) (out T, err error) {
	if err = c.Bind(requestData); err != nil {
		return out, &exception.ApplicationError{
			Code:    exception.ErrorCodeFailedBindingData,
			Message: err.Error(),
		}
	}

	//use the validator library to validate required fields
	if err = Validate(validate, requestData); err != nil {
		return out, err
	}

	return requestData, err
}

func Validate[T interface{}](validate validator.Validate, requestData T) error {
	if err := validate.Struct(requestData); err != nil {
		// Check if the error is a validation error
		var validationErrs validator.ValidationErrors
		if ok := errors.As(err, &validationErrs); ok {
			return &exception.ApplicationError{
				Code:    exception.ErrorCodeFailedBindingData,
				Message: "Validation failed",
				Details: getDetails(validationErrs),
			}
		}
	}
	return nil
}

func getDetails(validationErrs validator.ValidationErrors) (out []exception.ErrorDetail) {
	for _, vErr := range validationErrs {
		out = append(out, exception.ErrorDetail{
			Key:     vErr.Namespace(),
			Field:   vErr.Field(),
			Message: fmt.Sprintf("Failed on the '%s' tag", vErr.Tag()),
		})
	}

	return out
}
