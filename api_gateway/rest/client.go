package rest

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	jobOfferService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	messageService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	postService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceClientGrpc struct {
	AuthClient       authService.AuthServiceClient
	ProfileClient    profileService.ProfileServiceClient
	PostClient       postService.PostServiceClient
	ConnectionClient connectionService.ConnectionServiceClient
	JobOfferClient   jobOfferService.JobOfferServiceClient
	MessageClient    messageService.MessageServiceClient
}

func InitServiceClient(c *config.Config) *ServiceClientGrpc {
	client := &ServiceClientGrpc{
		AuthClient:       NewAuthClient(fmt.Sprintf("%s:%s", c.AuthHost, c.AuthPort)),
		ProfileClient:    NewProfileClient(fmt.Sprintf("%s:%s", c.ProfileHost, c.ProfilePort)),
		ConnectionClient: NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort)),
		PostClient:       NewPostClient(fmt.Sprintf("%s:%s", c.PostHost, c.PostPort)),
		JobOfferClient:   NewJobOfferClient(fmt.Sprintf("%s:%s", c.JobOfferHost, c.JobOfferPort)),
		MessageClient:    NewMessageClient(fmt.Sprintf("%s:%s", c.MessageHost, c.MessagePort)),
	}

	return client
}

func NewPostClient(address string) postService.PostServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return postService.NewPostServiceClient(conn)
}
func NewAuthClient(address string) authService.AuthServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return authService.NewAuthServiceClient(conn)
}

func NewProfileClient(address string) profileService.ProfileServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return profileService.NewProfileServiceClient(conn)
}

func NewConnectionClient(address string) connectionService.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return connectionService.NewConnectionServiceClient(conn)
}

func NewJobOfferClient(address string) jobOfferService.JobOfferServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return jobOfferService.NewJobOfferServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func NewMessageClient(address string) messageService.MessageServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return messageService.NewMessageServiceClient(conn)
}
