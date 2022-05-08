package error_mappers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InvalidArgumentErrorMapper struct{}

func NewInvalidArgumentErrorMapper() *InvalidArgumentErrorMapper {
	return &InvalidArgumentErrorMapper{}
}

func (m *InvalidArgumentErrorMapper) ToStatusError(domainErr error) (statusErr error) {
	return status.Error(codes.InvalidArgument, domainErr.Error())
}
