package errors

type AppError struct {
	Message string
	Code    string
	Status  int
	Details any
}

func (e *AppError) Error() string {
	return e.Message
}

func New(message string, status int) *AppError {
	return &AppError{
		Message: message,
		Status:  status,
		Code:    defaultCodeForStatus(status),
	}
}

func defaultCodeForStatus(status int) string {
	switch status {
	case 400:
		return ErrBadRequest
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	default:
		return ErrInternal
	}
}

func WithCode(err *AppError, code string) *AppError {
	err.Code = code
	return err
}

func WithDetails(err *AppError, details any) *AppError {
	err.Details = details
	return err
}
