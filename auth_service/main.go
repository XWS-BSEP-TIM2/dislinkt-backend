package main

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/startup"
	cfg "github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	fmt.Println("Auth service started.")
}
