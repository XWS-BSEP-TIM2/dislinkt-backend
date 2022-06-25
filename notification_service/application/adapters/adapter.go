package adapters

import (
	"crypto/tls"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func NewProfileClient(address string) profileService.ProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return profileService.NewProfileServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}

func NewConnectionClient(address string) connectionService.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return connectionService.NewConnectionServiceClient(conn)
}
