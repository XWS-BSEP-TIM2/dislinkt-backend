package profile

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/profile/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	config "github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func RegisterRoutes(r *gin.Engine, tracer opentracing.Tracer) {
	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	profileHandler := handler.InitProfileHandler(tracer)

	authorizedRoutes := r.Group("/profile")
	authorizedRoutes.Use(a.AuthRequired)
	authorizedRoutes.PUT("", a.Authorize("updateProfile", "update", false), profileHandler.Update)
	authorizedRoutes.PUT("/skills", a.Authorize("updateProfileSkills", "update", false), profileHandler.UpdateSkills)
	authorizedRoutes.POST("/changepassword", a.Authorize("changePassword", "update", false), profileHandler.ChangePassword)
	authorizedRoutes.GET("/admin-view", a.Authorize("getProfilesByAdmin", "read", false), profileHandler.Get)

	unauthorizedRoutes := r.Group("/profile")
	unauthorizedRoutes.GET("", profileHandler.Get)
	unauthorizedRoutes.GET("/:id", profileHandler.GetById)
	unauthorizedRoutes.POST("/search", profileHandler.Search)

}
