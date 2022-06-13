package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/application/adapters"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/startup/config"
)

type MessageService struct {
	store            persistence.MessageStore
	ConnectionClient connectionService.ConnectionServiceClient
	ProfileClient    profileService.ProfileServiceClient
}

func NewMessageService(store persistence.MessageStore, c *config.Config) *MessageService {
	return &MessageService{
		store:            store,
		ConnectionClient: adapters.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort)),
		ProfileClient:    adapters.NewProfileClient(fmt.Sprintf("%s:%s", c.ProfileHost, c.ProfilePort)),
	}
}

func (service *MessageService) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.MyContactsResponse, error) {
	//userID := request.UserID
	/*
		friends, err := service.ConnectionClient.GetFriends(ctx, &connectionService.GetRequest{UserID: userID})
		if err != nil {
			return nil, err
		}

	*/
	// for all msgID get contacts and make dto and return value

	return &pb.MyContactsResponse{Chats: []*pb.ChatPreview{{}}}, nil
}

func (service *MessageService) GetChat(ctx context.Context, request *pb.GetChatRequest) (*pb.ChatResponse, error) {

	return &pb.ChatResponse{Chat: nil}, nil
}
