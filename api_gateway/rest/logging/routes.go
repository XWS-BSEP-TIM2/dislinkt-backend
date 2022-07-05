package logging

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/logging/handler"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func RegisterRoutes(r *gin.Engine, tracer opentracing.Tracer) {

	configuration := config.NewConfig()
	a := security.NewAuthMiddleware(fmt.Sprintf("%s:%s", configuration.AuthHost, configuration.AuthPort))
	loggingHandler := handler.InitLoggingHandler(tracer)
	authorizedRoutes := r.Group("/events")

	authorizedRoutes.GET("/", a.Authorize("getAllEvents", "read", false), loggingHandler.GetAllEvents)

}
