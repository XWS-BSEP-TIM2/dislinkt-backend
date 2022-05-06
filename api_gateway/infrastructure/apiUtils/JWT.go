package apiUtils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

func GetJWTfromHTTPReq(r *http.Request) (string, error) {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		return "", status.Error(codes.Unauthenticated, `missing "Authorization" header`)
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
	}
	var token = strings.TrimPrefix(authHeader, prefix)
	return token, nil
}
