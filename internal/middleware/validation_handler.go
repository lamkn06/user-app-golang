package middleware

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/lamkn06/user-app-golang.git/pkg/exception"
)

func ParseValidationError(err error) *exception.ApplicationError {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var details []exception.ErrorDetail
		for _, vErr := range validationErrs {
			details = append(details, exception.ErrorDetail{
				Key:     vErr.Namespace(),
				Field:   vErr.Field(),
				Message: fmt.Sprintf("Failed on the '%s' tag", vErr.Tag()),
			})
		}

		return &exception.ApplicationError{
			Code:    exception.ErrorCodeValidation,
			Message: "Validation failed",
			Details: details,
		}
	}

	return &exception.ApplicationError{
		Code:    exception.ErrorCodeBadRequest,
		Message: err.Error(),
		Details: []exception.ErrorDetail{},
	}
}
