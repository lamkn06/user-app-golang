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

func ToApplicationError(err error, code string) *ApplicationError {
	appErr, ok := err.(*ApplicationError)
	if ok {
		return appErr
	}

	return &ApplicationError{
		Code:    code,
		Message: err.Error(),
		Details: []ErrorDetail{},
	}
}
