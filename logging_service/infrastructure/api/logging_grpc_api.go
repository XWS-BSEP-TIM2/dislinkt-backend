package api

import (
	//"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/application"
)

type LoggingHandler struct {
	pb.UnimplementedLoggingServiceServer
	service *application.LoggingService
}

func NewLoggingHandler(service *application.LoggingService) *LoggingHandler {
	return &LoggingHandler{
		service: service,
	}
}
