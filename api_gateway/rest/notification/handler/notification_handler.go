package handler

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type NotificationHandler struct {
	grpcClient *rest.ServiceClientGrpc
	tracer     opentracing.Tracer
}

func InitNotificationHandler(tracer opentracing.Tracer) *NotificationHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &NotificationHandler{grpcClient: client, tracer: tracer}
}

func (handler *NotificationHandler) GetAllNotifications(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetAllNotifications", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	notifications, err := notificationService.GetAllNotifications(ctx2, &pb.GetAllNotificationsRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &notifications)
}

func (handler *NotificationHandler) MarkAllAsSeen(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("MarkAllAsSeen", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	response, err := notificationService.MarkAllAsSeen(ctx2, &pb.MarkAllAsSeenRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}

func (handler *NotificationHandler) InsertNotification(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("InsertNotification", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	notificationService := handler.grpcClient.NotificationClient
	notificationDto := pb.Notification{}
	if err := ctx.BindJSON(&notificationDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	response, err := notificationService.InsertNotification(ctx2, &pb.InsertNotificationRequest{Notification: &notificationDto})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}

func (handler *NotificationHandler) GetUserSettings(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetUserSettings", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	response, err := notificationService.GetUserSettings(ctx2, &pb.GetUserSettingsRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}

func (handler *NotificationHandler) UpdateUserSettings(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("UpdateUserSettings", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	settingsCode := ctx.Param("code")
	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	response, err := notificationService.UpdateUserSettings(ctx2, &pb.UpdateUserSettingsRequest{UserID: dataFromToken.Id, SettingsCode: settingsCode})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}
