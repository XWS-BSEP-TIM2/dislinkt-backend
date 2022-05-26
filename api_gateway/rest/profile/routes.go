package profile

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/profile/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	config "github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	profileHandler := handler.InitProfileHandler()

	authorizedRoutes := r.Group("/profile")
	authorizedRoutes.Use(a.AuthRequired)
	authorizedRoutes.PUT("", a.Authorize("updateProfile", "update"), profileHandler.Update)
	authorizedRoutes.POST("/changepassword", a.Authorize("changePassword", "update"), profileHandler.ChangePassword)

	unauthorizedRoutes := r.Group("/profile")
	unauthorizedRoutes.GET("", profileHandler.Get)
	unauthorizedRoutes.GET("/:id", profileHandler.GetById)
	unauthorizedRoutes.POST("/search", profileHandler.Search)

}
