package post

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/post/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	postHandler := handler.InitPostHandler()
	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))

	authorizedRoutes := r.Group("/posts")
	authorizedRoutes.POST("", a.Authorize("createPost", "create"), postHandler.CreatePost)
	authorizedRoutes.POST("/:id/comments", a.Authorize("createComment", "create"), postHandler.CreateComment)
	authorizedRoutes.POST("/:id/likes", a.Authorize("likePost", "create"), postHandler.LikePost)
	authorizedRoutes.DELETE(":id/likes/:like-id", a.Authorize("removeLike", "delete"), postHandler.RemoveLike)
	authorizedRoutes.POST("/:id/dislikes", a.Authorize("dislikePost", "create"), postHandler.DislikePost)
	authorizedRoutes.DELETE(":id/dislikes/:dislike-id", a.Authorize("removeDislike", "delete"), postHandler.RemoveDislike)

	unauthorizedRoutes := r.Group("/posts")
	unauthorizedRoutes.GET("", postHandler.Get)
	unauthorizedRoutes.GET("/:id", postHandler.GetPostById)
	unauthorizedRoutes.GET("/:id/comments", postHandler.GetPostComments)
	unauthorizedRoutes.GET("/:id/comments/:comment-id", postHandler.GetPostComment)
	unauthorizedRoutes.GET("/:id/likes", postHandler.GetLikes)
	unauthorizedRoutes.GET("/:id/likes/:like-id", postHandler.GetLike)
	unauthorizedRoutes.GET("/:id/dislikes", postHandler.GetDislikes)
	unauthorizedRoutes.GET("/:id/dislikes/:dislike-id", postHandler.GetDislike)

}
