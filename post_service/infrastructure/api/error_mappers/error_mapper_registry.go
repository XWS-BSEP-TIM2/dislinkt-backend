package error_mappers

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
)

type ErrorMapperRegistry struct {
}

func NewErrorMapperRegistry() *ErrorMapperRegistry {
	return &ErrorMapperRegistry{}
}

func (registry *ErrorMapperRegistry) ToStatusError(domainErr error) (statusErr error) {
	switch domainErr.(type) {
	case *errors.EntityNotFoundError:
		return NewEntityNotFoundErrorMapper().ToStatusError(domainErr)
	case *errors.EntityForbiddenError:
		return NewEntityForbiddenErrorMapper().ToStatusError(domainErr)
	default:
		return NewDefaultErrorMapper().ToStatusError(domainErr)
	}

}
