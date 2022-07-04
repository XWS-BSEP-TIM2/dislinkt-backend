package notification

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/notification/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func RegisterRoutes(r *gin.Engine, tracer opentracing.Tracer) {

	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	notificationHandler := handler.InitNotificationHandler(tracer)
	authorizedRoutes := r.Group("/notifications")

	authorizedRoutes.GET("/:userId", a.Authorize("getAllNotifications", "read", false), notificationHandler.GetAllNotifications)
	authorizedRoutes.POST("/", a.Authorize("insertNotification", "create", false), notificationHandler.InsertNotification)
	authorizedRoutes.PUT("/:userId", a.Authorize("markAllNotificationsAsRead", "update", false), notificationHandler.MarkAllAsSeen)
	authorizedRoutes.GET("/settings/:userId", a.Authorize("getUserSettings", "read", false), notificationHandler.GetUserSettings)
	authorizedRoutes.PUT("/settings/:userId/:code", a.Authorize("updateUserSettings", "update", false), notificationHandler.UpdateUserSettings)
}
