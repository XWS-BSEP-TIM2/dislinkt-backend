package main

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup"
	cfg "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	fmt.Println("Post service started.")
}
