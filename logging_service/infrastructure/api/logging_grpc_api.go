package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "LoggInfo")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.LoggInfo(ctx2, request)
}

func (handler *LoggingHandler) LoggError(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggError")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.LoggError(ctx2, request)
}

func (handler *LoggingHandler) LoggWarning(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggWarning")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.LoggWarning(ctx2, request)
}

func (handler *LoggingHandler) LoggSuccess(ctx context.Context, request *pb.LogRequest) (*pb.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "LoggSuccess")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.service.LoggSuccess(ctx2, request)
}
