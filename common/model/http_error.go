package model

type HttpError struct {
	code    int
	message string
}

func (err *HttpError) Error() string {
	return err.message
}

func NewHttpError(code int, message string) *HttpError {
	return &HttpError{code: code, message: message}
}

func (err *HttpError) GetCode() int {
	return err.code
}

func (err *HttpError) GetMessage() string {
	return err.message
}
