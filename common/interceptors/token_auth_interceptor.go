package interceptors

import (
	"context"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

var nonAuthMethods map[string]bool = map[string]bool{
	"/connection_service.ConnectionService/Register": true,
	"/auth_service.AuthService/Register":             true,
	"/auth_service.AuthService/Login":                true,
	"/profile_service.ProfileService/CreateProfile":  true,
	"/profile_service.ProfileService/Get":            true,
	"/post_service.PostService/GetPost":              true,
	"/post_service.PostService/CreatePost":           true,
	"/post_service.PostService/GetPosts":             true,
	"/post_service.PostService/CreateComment":        true,
	"/post_service.PostService/GetComment":           true,
	"/post_service.PostService/GetComments":          true,
}

func TokenAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	method, _ := grpc.Method(ctx)
	if nonAuthMethods[method] == true {
		return handler(ctx, req)
	}

	auth, err := extractHeader(ctx, "authorization")
	if err != nil {
		return ctx, err
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
	}
	authClient := NewAuthClient("localhost:8081")
	var token = strings.TrimPrefix(auth, prefix)
	result, err := authClient.Validate(context.TODO(), &authService.ValidateRequest{Token: token})
	if result.Status != 200 {
		return ctx, status.Error(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
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
