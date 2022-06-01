package helper

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func ExtractTokenFromContext(ctx context.Context) (string, error) {
	auth, err := extractHeader(ctx, "authorization")
	if err != nil {
		return "", err
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return "invalid_auth_header", status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
	}
	var token = strings.TrimPrefix(auth, prefix)
	return token, nil
}

func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}
