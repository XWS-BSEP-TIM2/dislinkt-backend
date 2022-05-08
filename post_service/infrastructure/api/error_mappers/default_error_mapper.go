package error_mappers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DefaultErrorMapper struct{}

func NewDefaultErrorMapper() *DefaultErrorMapper {
	return &DefaultErrorMapper{}
}

func (m *DefaultErrorMapper) ToStatusError(domainErr error) (statusErr error) {
	return status.Error(codes.Aborted, domainErr.Error())
}
