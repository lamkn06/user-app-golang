package exception

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

	if code == ErrorCodeInternalServerError {
		return &ApplicationError{
			Code: code,
		}
	}

	return &ApplicationError{
		Code:    code,
		Message: err.Error(),
		Details: []ErrorDetail{},
	}
}
