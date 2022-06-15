package api

import (
	"context"
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

func (handler *LoggingHandler) LoggInfo(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	return handler.service.LoggInfo(ctx, request)
}

func (handler *LoggingHandler) LoggError(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	return handler.service.LoggError(ctx, request)
}

func (handler *LoggingHandler) LoggWarning(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	return handler.service.LoggWarning(ctx, request)
}

func (handler *LoggingHandler) LoggSuccess(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	return handler.service.LoggSuccess(ctx, request)
}
