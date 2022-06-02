package main

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/auth"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/connection"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/job_offer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/post"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/profile"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	"github.com/gin-gonic/gin"
)

func main() {
	//config := config.NewConfig()
	//server := startup.NewServer(config)
	//server.Start()

	c := config.NewConfig()

	r := gin.Default()
	r.Use(security.CORSMiddleware())

	auth.RegisterRoutes(r)
	profile.RegisterRoutes(r)
	post.RegisterRoutes(r)
	connection.RegisterRoutes(r)
	job_offer.RegisterRoutes(r)

	//r.Run(":" + c.Port)

	err := r.RunTLS(":"+c.Port, "./certificates/dislinkt_gateway.crt", "./certificates/dislinkt_gateway.key")
	if err != nil {
		return
	}
}
