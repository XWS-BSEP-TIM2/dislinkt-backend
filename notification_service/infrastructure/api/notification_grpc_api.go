package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	service *application.NotificationService
}

func NewNotificationHandler(service *application.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: service,
	}
}

func (handler *NotificationHandler) GetAllNotifications(ctx context.Context, request *pb.GetAllNotificationsRequest) (*pb.GetAllNotificationsResponse, error) {
	return handler.service.GetAllNotifications(ctx, request)
}

func (handler *NotificationHandler) MarkAllAsSeen(ctx context.Context, request *pb.MarkAllAsSeenRequest) (*pb.MarkAllAsSeenResponse, error) {
	return handler.service.MarkAllAsSeen(ctx, request)
}

func (handler *NotificationHandler) InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest) (*pb.InsertNotificationRequestResponse, error) {
	return handler.service.InsertNotification(ctx, request)
}

func (handler *NotificationHandler) GetUserSettings(ctx context.Context, request *pb.GetUserSettingsRequest) (*pb.GetUserSettingsResponse, error) {
	return handler.service.GetUserSettings(ctx, request)
}

func (handler *NotificationHandler) UpdateUserSettings(ctx context.Context, request *pb.UpdateUserSettingsRequest) (*pb.GetUserSettingsResponse, error) {
	return handler.service.UpdateUserSettings(ctx, request)
}
