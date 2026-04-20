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
