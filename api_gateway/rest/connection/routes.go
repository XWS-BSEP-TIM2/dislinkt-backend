package connection

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/connection/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func RegisterRoutes(r *gin.Engine, tracer opentracing.Tracer) {

	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	connectionHandler := handler.InitConnectionHandler(tracer)
	authorizedRoutes := r.Group("/connection")
	authorizedRoutes.GET("/friends/:id", a.Authorize("getFriends", "read", false), connectionHandler.GetFriends)
	authorizedRoutes.GET("/blocks", a.Authorize("getBlockedUsers", "read", false), connectionHandler.GetBlocks)
	authorizedRoutes.GET("/friends-requests", a.Authorize("getFriendRequests", "read", false), connectionHandler.GetFriendRequests)
	authorizedRoutes.POST("/friend", a.Authorize("addFriend", "create", false), connectionHandler.AddFriend)
	authorizedRoutes.POST("/remove-friend", a.Authorize("removeFriend", "delete", false), connectionHandler.RemoveFriend)
	authorizedRoutes.POST("/block", a.Authorize("blockUser", "create", false), connectionHandler.BlockUser)
	authorizedRoutes.POST("/unblock", a.Authorize("unblockUser", "read", false), connectionHandler.UnblockUser)
	authorizedRoutes.POST("/friend-request", a.Authorize("createFriendRequest", "create", false), connectionHandler.CreateFriendRequest)
	authorizedRoutes.POST("/remove-friend-request", a.Authorize("removeFriendRequest", "delete", false), connectionHandler.DeleteFriendRequest)
	authorizedRoutes.GET("/recommendation", a.Authorize("getRecommendation", "read", false), connectionHandler.GetRecommendation)
	authorizedRoutes.GET("/:ida/detail/:idb", a.Authorize("getConnectionDetails", "read", false), connectionHandler.GetConnectionDetails)
}
