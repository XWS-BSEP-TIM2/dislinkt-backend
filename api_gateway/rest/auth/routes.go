package auth

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/auth/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	authHandler := handler.InitAuthHandler()
	routes := r.Group("")
	routes.POST("/api/register", authHandler.Register) //TODO: "/register" treba na frontu obrisati api, priveremo
	routes.POST("/login", authHandler.Login)
	routes.GET("/verify/:username/:code", authHandler.Verify)
}
