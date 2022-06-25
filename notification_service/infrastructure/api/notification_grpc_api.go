package api

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application"
)

type NotificationHandler struct {
	pb.UnimplementedMessageServiceServer
	service *application.NotificationService
}

func NewMessageHandler(service *application.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}
