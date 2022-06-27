package handler

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type NotificationHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func InitNotificationHandler() *NotificationHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &NotificationHandler{grpcClient: client}
}

func (handler *NotificationHandler) GetAllNotifications(ctx *gin.Context) {
	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	notifications, err := notificationService.GetAllNotifications(ctx, &pb.GetAllNotificationsRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &notifications)
}

func (handler *NotificationHandler) MarkAllAsSeen(ctx *gin.Context) {
	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	response, err := notificationService.MarkAllAsSeen(ctx, &pb.MarkAllAsSeenRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}

func (handler *NotificationHandler) InsertNotification(ctx *gin.Context) {
	notificationService := handler.grpcClient.NotificationClient
	notificationDto := pb.Notification{}
	if err := ctx.BindJSON(&notificationDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	response, err := notificationService.InsertNotification(ctx, &pb.InsertNotificationRequest{Notification: &notificationDto})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}

func (handler *NotificationHandler) GetUserSettings(ctx *gin.Context) {
	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	response, err := notificationService.GetUserSettings(ctx, &pb.GetUserSettingsRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}

func (handler *NotificationHandler) UpdateUserSettings(ctx *gin.Context) {
	settingsCode := ctx.Param("code")
	notificationService := handler.grpcClient.NotificationClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	response, err := notificationService.UpdateUserSettings(ctx, &pb.UpdateUserSettingsRequest{UserID: dataFromToken.Id, SettingsCode: settingsCode})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &response)
}
