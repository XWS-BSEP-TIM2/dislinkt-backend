package main

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup"
	cfg "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
