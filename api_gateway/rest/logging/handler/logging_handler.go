package handler

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type LoggingHandler struct {
	grpcClient *rest.ServiceClientGrpc
	tracer     opentracing.Tracer
}

func InitLoggingHandler(tracer opentracing.Tracer) *LoggingHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &LoggingHandler{grpcClient: client, tracer: tracer}
}

func (handler *LoggingHandler) GetAllEvents(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetAllEvents", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	loggingService := handler.grpcClient.LoggingClient
	events, err := loggingService.GetEvents(ctx2, &pb.Empty{})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &events)
}
