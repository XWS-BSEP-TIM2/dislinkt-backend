package job_offer

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/job_offer/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	jobOfferHandler := handler.InitJobOfferHandler()

	authorizedRoutes := r.Group("/job-offer")
	authorizedRoutes.Use(a.AuthRequired)

	authorizedRoutes.PUT("", a.Authorize("updateJobOffer", "update", false), jobOfferHandler.Update)
	unauthorizedRoutes := r.Group("/job-offer")
	unauthorizedRoutes.GET("", jobOfferHandler.Get)
	unauthorizedRoutes.GET("/:id", jobOfferHandler.GetById)
	unauthorizedRoutes.POST("/search", jobOfferHandler.Search)
	unauthorizedRoutes.POST("", jobOfferHandler.Create)
	unauthorizedRoutes.GET("/user-offers/:id", jobOfferHandler.GetUserJobOffers)
	unauthorizedRoutes.DELETE("/:id", jobOfferHandler.DeleteOffer)
}
