package error_mappers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EntityForbiddenErrorMapper struct{}

func NewEntityForbiddenErrorMapper() *EntityForbiddenErrorMapper {
	return &EntityForbiddenErrorMapper{}
}

func (m *EntityForbiddenErrorMapper) ToStatusError(domainErr error) (statusErr error) {
	return status.Error(codes.PermissionDenied, domainErr.Error())
}
