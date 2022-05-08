package service_clients

import (
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"log"
)

func NewProfileClient(address string) profileService.ProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to [Profile service]: %v", err)
	}
	return profileService.NewProfileServiceClient(conn)
}
