package connection

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/connection/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	connectionHandler := handler.InitConnectionHandler()
	authorizedRoutes := r.Group("/connection")
	authorizedRoutes.GET("/friends/user/:id", a.Authorize("getFriends", "read"), connectionHandler.GetFriends)
	authorizedRoutes.GET("/user/:id/blockeds", a.Authorize("getBlockedUsers", "read"), connectionHandler.GetBlocks)
	authorizedRoutes.GET("/friend-request", a.Authorize("getFriendRequests", "read"), connectionHandler.GetFriendRequests)
	authorizedRoutes.POST("/friend", a.Authorize("addFriend", "create"), connectionHandler.AddFriend)
	authorizedRoutes.POST("/remove-friend", a.Authorize("removeFriend", "delete"), connectionHandler.RemoveFriend)
	authorizedRoutes.POST("/block", a.Authorize("blockUser", "create"), connectionHandler.BlockUser)
	authorizedRoutes.POST("/unblock", a.Authorize("unblockUser", "read"), connectionHandler.UnblockUser)
	authorizedRoutes.POST("/friend-request", a.Authorize("createFriendRequest", "create"), connectionHandler.CreateFriendRequest)
	authorizedRoutes.POST("/remove-friend-request", a.Authorize("removeFriendRequest", "delete"), connectionHandler.DeleteFriendRequest)
	authorizedRoutes.GET("/recommendation", a.Authorize("getRecommendation", "read"), connectionHandler.GetRecommendation)
	authorizedRoutes.GET("/:ida/detail/idb", a.Authorize("getConnectionDetails", "read"), connectionHandler.GetConnectionDetails)
}
