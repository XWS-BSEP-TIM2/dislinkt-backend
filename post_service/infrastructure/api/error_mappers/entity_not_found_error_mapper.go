package error_mappers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EntityNotFoundErrorMapper struct{}

func NewEntityNotFoundErrorMapper() *EntityNotFoundErrorMapper {
	return &EntityNotFoundErrorMapper{}
}

func (m *EntityNotFoundErrorMapper) ToStatusError(domainErr error) (statusErr error) {
	return status.Error(codes.NotFound, domainErr.Error())
}
