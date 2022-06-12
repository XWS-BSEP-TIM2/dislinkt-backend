package application

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/persistence"
)

type MessageService struct {
	store persistence.MessageStore
}

func NewMessageService(store persistence.MessageStore) *MessageService {
	return &MessageService{
		store: store,
	}
}

func (service *MessageService) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.MyContactsResponse, error) {
	// call connection service get my friends and msgIDs
	// for all msgID get contacts and make dto and return value
	return &pb.MyContactsResponse{Chats: []*pb.ChatPreview{{}}}, nil
}

/*
   func (service *MessageService) Get(ctx context.Context, id primitive.ObjectID) (*domain.Chat, error) {
   	return service.store.Get(ctx, id)
   }

   func (service *MessageService) GetAll(ctx context.Context) ([]*domain.Chat, error) {
   	return service.store.GetAll(ctx)
   }

   func (service *MessageService) Insert(ctx context.Context, chat *domain.Chat) {
   	service.store.Insert(ctx, chat)

   }

   func (service *MessageService) Update(ctx context.Context, chat *domain.Chat) {
   	service.store.Update(ctx, chat)
   }

   func (service *MessageService) Search(ctx context.Context, search string) ([]*domain.Chat, error) {
   	return service.store.Search(ctx, search)
   }

*/
