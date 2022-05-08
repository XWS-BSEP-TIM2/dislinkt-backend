package errors

// EntityNotFoundError
// throw when there is no entity that matches some description
type EntityNotFoundError struct {
	message string
}

func NewEntityNotFoundError(message string) *EntityNotFoundError {
	return &EntityNotFoundError{message}
}

func (err *EntityNotFoundError) Error() string { return err.message }
