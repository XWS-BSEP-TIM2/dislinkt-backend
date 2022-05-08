package errors

type InvalidArgumentError struct {
	message string
}

func NewInvalidArgumentError(message string) *InvalidArgumentError {
	return &InvalidArgumentError{message}
}

func (err *InvalidArgumentError) Error() string { return err.message }
