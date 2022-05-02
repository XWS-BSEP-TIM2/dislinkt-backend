package main

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/startup"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
