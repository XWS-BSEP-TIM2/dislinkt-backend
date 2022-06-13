package messages

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/messages/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	messageHandler := handler.InitMessageHandler()
	authorizedRoutes := r.Group("/messages")

	authorizedRoutes.GET("/contacts", a.Authorize("getContacts", "read", false), messageHandler.GetMyContacts)
	authorizedRoutes.GET("/chat/:chatId", a.Authorize("getChat", "read", false), messageHandler.GetChat)
	authorizedRoutes.POST("/chat/send", a.Authorize("sendMessage", "update", false), messageHandler.SendMessage)
	authorizedRoutes.POST("/chat/seen", a.Authorize("setSeen", "update", false), messageHandler.SetSeen)

}
