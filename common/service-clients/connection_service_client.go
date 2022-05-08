package service_clients

import (
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"log"
)

func NewConnectionClient(address string) connectionService.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to [Connection service]: %v", err)
	}
	return connectionService.NewConnectionServiceClient(conn)
}
