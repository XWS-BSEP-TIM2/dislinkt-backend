package service_clients

import (
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"log"
)

func NewAuthClient(address string) authService.AuthServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to [Auth service]: %v", err)
	}
	return authService.NewAuthServiceClient(conn)
}
