package main

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/auth"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/connection"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/job_offer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/messages"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/notification"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/post"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/profile"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {
	//config := config.NewConfig()
	//server := startup.NewServer(config)
	//server.Start()

	c := config.NewConfig()

	r := gin.Default()
	r.Use(security.CORSMiddleware())

	// metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(r)

	// tracing
	tracer, _ := tracer.Init("api_gateway")
	opentracing.SetGlobalTracer(tracer)

	auth.RegisterRoutes(r, tracer)
	profile.RegisterRoutes(r, tracer)
	post.RegisterRoutes(r, tracer)
	connection.RegisterRoutes(r, tracer)
	job_offer.RegisterRoutes(r)
	messages.RegisterRoutes(r)
	notification.RegisterRoutes(r)

	//r.Run(":" + c.Port)

	err := r.RunTLS(":"+c.Port, "./certificates/dislinkt_gateway.crt", "./certificates/dislinkt_gateway.key")
	if err != nil {
		return
	}

}
