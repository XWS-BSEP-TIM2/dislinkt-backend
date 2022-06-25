package application

import (
	"crypto/tls"
	"fmt"
	notificationService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}

func NewNotificationClient(address string) notificationService.NotificationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return notificationService.NewNotificationServiceClient(conn)
}

func NewProfileClient(address string) profileService.ProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return profileService.NewProfileServiceClient(conn)
}
