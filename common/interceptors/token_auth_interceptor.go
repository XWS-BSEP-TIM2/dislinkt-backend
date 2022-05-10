package interceptors

import (
	"context"
	"fmt"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
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
	"/post_service.PostService/GiveLike":             true,
	"/post_service.PostService/GetLike":              true,
	"/post_service.PostService/GetLikes":             true,
	"/post_service.PostService/UndoLike":             true,
	"/post_service.PostService/GiveDislike":          true,
	"/post_service.PostService/GetDislike":           true,
	"/post_service.PostService/GetDislikes":          true,
	"/post_service.PostService/UndoDislike":          true,
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

	authEndpoint := fmt.Sprintf("%s:%s", goDotEnvVariable("AUTH_SERVICE_HOST"), goDotEnvVariable("AUTH_SERVICE_PORT"))
	authClient := NewAuthClient(authEndpoint)
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

func goDotEnvVariable(key string) string {
	godotenv.Load("../.env")

	/*
		if err != nil {
			log.Fatalf("Error loading .env file")
		} */
	return os.Getenv(key)
}
