package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "GetAllNotifications")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.GetAllNotifications(ctx2, request)
}

func (handler *NotificationHandler) MarkAllAsSeen(ctx context.Context, request *pb.MarkAllAsSeenRequest) (*pb.MarkAllAsSeenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "MarkAllAsSeen")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.MarkAllAsSeen(ctx2, request)
}

func (handler *NotificationHandler) InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest) (*pb.InsertNotificationRequestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertNotification")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.InsertNotification(ctx2, request)
}

func (handler *NotificationHandler) GetUserSettings(ctx context.Context, request *pb.GetUserSettingsRequest) (*pb.GetUserSettingsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUserSettings")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.GetUserSettings(ctx2, request)
}

func (handler *NotificationHandler) UpdateUserSettings(ctx context.Context, request *pb.UpdateUserSettingsRequest) (*pb.GetUserSettingsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSettings")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.UpdateUserSettings(ctx2, request)
}
