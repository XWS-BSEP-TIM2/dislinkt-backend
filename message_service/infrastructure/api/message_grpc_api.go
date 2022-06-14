package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/application"
)

type MessageHandler struct {
	pb.UnimplementedMessageServiceServer
	service *application.MessageService
}

func NewMessageHandler(service *application.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

func (handler *MessageHandler) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.MyContactsResponse, error) {
	return handler.service.GetMyContacts(ctx, request)
}

func (handler *MessageHandler) GetChat(ctx context.Context, request *pb.GetChatRequest) (*pb.ChatResponse, error) {
	return handler.service.GetChat(ctx, request)
}

func (handler *MessageHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.ActionResult, error) {
	return handler.service.SendMessage(ctx, request)
}

func (handler *MessageHandler) SetSeen(ctx context.Context, request *pb.SetSeenRequest) (*pb.ActionResult, error) {
	return handler.service.SetSeen(ctx, request)
}
