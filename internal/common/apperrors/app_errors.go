package apperrors

type AppError struct {
	Err        error
	Type       string
	ClientMsg  string
	StatusCode int
	RenderType RenderType
}

func (e AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, errorType string, clientMsg string, statusCode int, renderType RenderType) *AppError {
	return &AppError{
		Err:        err,
		Type:       errorType,
		ClientMsg:  clientMsg,
		StatusCode: statusCode,
		RenderType: renderType,
	}
}
