package exception

import (
	"strings"
)

type ErrorDetail struct {
	Key     string `json:"key"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ApplicationError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details"`
}

func (e *ApplicationError) Error() string {
	return e.Message
}

func (e *ApplicationError) HTTPStatus() int {
	switch e.Code {
	case ErrorCodeValidation:
		return 400
	case ErrorCodeBadRequest:
		return 400
	case ErrorCodeUnauthorized:
		return 401
	case ErrorCodeForbidden:
		return 403
	case ErrorCodeNotFound:
		return 404
	case ErrorCodeTooManyRequests:
		return 429
	default:
		return 500
	}
}

func ToApplicationError(err error, code string) *ApplicationError {
	appErr, ok := err.(*ApplicationError)
	if ok {
		return appErr
	}

	// Hide database-related errors for security
	if code == ErrorCodeInternalServerError || isDatabaseError(err) {
		return &ApplicationError{
			Code:    code,
			Message: "Internal server error",
			Details: []ErrorDetail{},
		}
	}

	return &ApplicationError{
		Code:    code,
		Message: err.Error(),
		Details: []ErrorDetail{},
	}
}

// isDatabaseError checks if the error is related to database operations
func isDatabaseError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())

	// Check for common database error patterns
	databaseKeywords := []string{
		"database",
		"connection",
		"sql",
		"postgres",
		"mysql",
		"constraint",
		"foreign key",
		"duplicate",
		"unique",
		"not null",
		"timeout",
		"deadlock",
		"rollback",
		"transaction",
	}

	for _, keyword := range databaseKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}

	return false
}
