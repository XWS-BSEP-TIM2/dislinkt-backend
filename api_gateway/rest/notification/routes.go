package notification

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/notification/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	notificationHandler := handler.InitNotificationHandler()
	authorizedRoutes := r.Group("/notifications")

	authorizedRoutes.GET("/:userId", a.Authorize("getAll", "read", false), notificationHandler.GetAllNotifications)
	authorizedRoutes.POST("/", a.Authorize("insertNotification", "write", false), notificationHandler.InsertNotification)
	authorizedRoutes.PUT("/:userId", a.Authorize("markAllAsRead", "update", false), notificationHandler.MarkAllAsSeen)

}
