package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	jobOfferService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	loggingService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	messageService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	notificationService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	postService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

type ServiceClientGrpc struct {
	AuthClient         authService.AuthServiceClient
	ProfileClient      profileService.ProfileServiceClient
	PostClient         postService.PostServiceClient
	ConnectionClient   connectionService.ConnectionServiceClient
	JobOfferClient     jobOfferService.JobOfferServiceClient
	MessageClient      messageService.MessageServiceClient
	LoggingClient      loggingService.LoggingServiceClient
	NotificationClient notificationService.NotificationServiceClient
}

func InitServiceClient(c *config.Config) *ServiceClientGrpc {
	client := &ServiceClientGrpc{
		AuthClient:         NewAuthClient(fmt.Sprintf("%s:%s", c.AuthHost, c.AuthPort)),
		ProfileClient:      NewProfileClient(fmt.Sprintf("%s:%s", c.ProfileHost, c.ProfilePort)),
		ConnectionClient:   NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort)),
		PostClient:         NewPostClient(fmt.Sprintf("%s:%s", c.PostHost, c.PostPort)),
		JobOfferClient:     NewJobOfferClient(fmt.Sprintf("%s:%s", c.JobOfferHost, c.JobOfferPort)),
		MessageClient:      NewMessageClient(fmt.Sprintf("%s:%s", c.MessageHost, c.MessagePort)),
		LoggingClient:      NewLoggingClient(fmt.Sprintf("%s:%s", c.LoggingHost, c.LoggingPort)),
		NotificationClient: NewNotificationClient(fmt.Sprintf("%s:%s", c.NotificationHost, c.NotificationPort)),
	}

	return client
}

func NewPostClient(address string) postService.PostServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return postService.NewPostServiceClient(conn)
}
func NewAuthClient(address string) authService.AuthServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return authService.NewAuthServiceClient(conn)
}

func NewProfileClient(address string) profileService.ProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return profileService.NewProfileServiceClient(conn)
}

func NewConnectionClient(address string) connectionService.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return connectionService.NewConnectionServiceClient(conn)
}

func NewJobOfferClient(address string) jobOfferService.JobOfferServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return jobOfferService.NewJobOfferServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.DialContext(
		context.Background(),
		address,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		// tracer
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		),
		grpc.WithStreamInterceptor(
			grpc_opentracing.StreamClientInterceptor(
				grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
			),
		))

	//return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}

func NewMessageClient(address string) messageService.MessageServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Message service: %v", err)
	}
	return messageService.NewMessageServiceClient(conn)
}

func NewNotificationClient(address string) notificationService.NotificationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Notification service: %v", err)
	}
	return notificationService.NewNotificationServiceClient(conn)
}

func NewLoggingClient(address string) loggingService.LoggingServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Logging service: %v", err)
	}
	return loggingService.NewLoggingServiceClient(conn)
}
