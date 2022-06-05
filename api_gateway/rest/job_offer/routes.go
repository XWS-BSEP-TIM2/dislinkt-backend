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
	authorizedRoutes.GET("/:id", a.Authorize("getJobOffer", "read", false), jobOfferHandler.GetById)
	authorizedRoutes.POST("/search", a.Authorize("searchJobOffers", "read", false), jobOfferHandler.Search)
	authorizedRoutes.POST("", a.Authorize("createJobOffer", "create", false), jobOfferHandler.Create)
	authorizedRoutes.GET("/user-offers/:id", a.Authorize("getUserJobOffers", "read", false), jobOfferHandler.GetUserJobOffers)
	authorizedRoutes.DELETE("/:id", a.Authorize("deleteJobOffer", "delete", false), jobOfferHandler.DeleteOffer)

	unauthorizedRoutes := r.Group("/job-offer")
	unauthorizedRoutes.GET("", jobOfferHandler.Get)

	apiRoutes := r.Group("api/job-offer")
	apiRoutes.PUT("", a.Authorize("updateJobOffer", "update", true), jobOfferHandler.Update)
	apiRoutes.POST("", a.Authorize("createJobOffer", "create", true), jobOfferHandler.CreateFromExternalApp)
	apiRoutes.DELETE("/:id", a.Authorize("deleteJobOffer", "delete", true), jobOfferHandler.DeleteOffer)
}
