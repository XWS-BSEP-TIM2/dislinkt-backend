package main

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
