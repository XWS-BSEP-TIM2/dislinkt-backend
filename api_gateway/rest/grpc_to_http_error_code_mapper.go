package rest

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

type GrpcToHttpErrorCodeMapper struct {
	m map[codes.Code]int
}

func NewGrpcToHttpErrorCodeMapper() *GrpcToHttpErrorCodeMapper {
	return &GrpcToHttpErrorCodeMapper{
		m: map[codes.Code]int{
			codes.NotFound:         http.StatusNotFound,
			codes.InvalidArgument:  http.StatusBadRequest,
			codes.PermissionDenied: http.StatusForbidden,
			codes.Aborted:          http.StatusInternalServerError,
		},
	}
}

func (mapper *GrpcToHttpErrorCodeMapper) MapGrpcToHttpError(code codes.Code) (httpStatus int) {
	return mapper.m[code]
}
