package logging_service_adapter

import (
	"context"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters"
	"google.golang.org/grpc/peer"
)

type LoggingServiceAdapter struct {
	address string
}

func NewLoggingServiceAdapter(address string) *LoggingServiceAdapter {
	return &LoggingServiceAdapter{address: address}
}

func (l *LoggingServiceAdapter) Log(ctx context.Context, logType, serviceFunctionName, userID, description string) {
	span := tracer.StartSpanFromContext(ctx, "Log")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	logService := adapters.NewLoggingClient(l.address)
	ipAddress := ""
	p, ok := peer.FromContext(ctx)
	if ok {
		ipAddress = p.Addr.String()
	}
	if logType == "ERROR" {
		logService.LoggError(ctx2, &loggingS.LogRequest{ServiceName: "POST_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "SUCCESS" {
		logService.LoggSuccess(ctx2, &loggingS.LogRequest{ServiceName: "POST_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "WARNING" {
		logService.LoggWarning(ctx2, &loggingS.LogRequest{ServiceName: "POST_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "INFO" {
		logService.LoggInfo(ctx2, &loggingS.LogRequest{ServiceName: "POST_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	}
}
