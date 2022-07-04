package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "GetMyContacts")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.GetMyContacts(ctx2, request)
}

func (handler *MessageHandler) GetChat(ctx context.Context, request *pb.GetChatRequest) (*pb.ChatResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetChat")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.GetChat(ctx2, request)
}

func (handler *MessageHandler) SendMessage(ctx context.Context, request *pb.SendMessageRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "SendMessage")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.SendMessage(ctx2, request)
}

func (handler *MessageHandler) SetSeen(ctx context.Context, request *pb.SetSeenRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "SetSeen")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.SetSeen(ctx2, request)
}

func (handler *MessageHandler) CreateChat(ctx context.Context, request *pb.CreateChatRequest) (*pb.CreateChatResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateChat")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.CreateChat(ctx2, request)
}
