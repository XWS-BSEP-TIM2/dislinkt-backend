package adapters

import (
	"crypto/tls"
	"fmt"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	loggingService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	notificationService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func NewAuthClient(address string) authService.AuthServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Auth service: %v", err)
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

func NewNotificationClient(address string) notificationService.NotificationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
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

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}
