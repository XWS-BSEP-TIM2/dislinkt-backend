package auth

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/auth/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func RegisterRoutes(r *gin.Engine, tracer opentracing.Tracer) {
	authHandler := handler.InitAuthHandler(tracer)
	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	routes := r.Group("")
	routes.POST("/api/register", authHandler.Register) //TODO: "/register" treba na frontu obrisati api, priveremo
	routes.POST("/login", authHandler.Login)
	routes.GET("/verify/:username/:code", authHandler.Verify)
	routes.GET("/login/verify/:username", authHandler.ResendVerify)
	routes.GET("/login/recovery/:username", authHandler.GetRecovery)
	routes.POST("/login/recovery", authHandler.Recover)
	routes.POST("magic-link-login/send-mail", authHandler.SendMailForMagicLinkRegistration)
	routes.POST("magic-link-login", authHandler.MagicLinkLogin)
	routes.GET("/api-token/:userId", authHandler.GenerateApiToken)
	routes.GET("/test", a.Authorize("test", "create", true), authHandler.Test)
	routes.GET("/two-factor/:id", authHandler.GenerateQrCode)
	routes.POST("/two-factor", authHandler.Verify2Factor)
}
