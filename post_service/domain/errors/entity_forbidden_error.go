package errors

// EntityForbiddenError
// throw when there is entity that matches some description but requester doesn't have authorization
type EntityForbiddenError struct {
	message string
}

func NewEntityForbiddenError(message string) *EntityForbiddenError {
	return &EntityForbiddenError{message}
}

func (err *EntityForbiddenError) Error() string { return err.message }
