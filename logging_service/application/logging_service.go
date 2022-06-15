package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/persistence"
)

const (
	ERROR   = "ERROR"
	INFO    = "INFO"
	WARNING = "WARNING"
	SUCCESS = "SUCCESS"
)

type LoggingService struct {
	store persistence.LoggingStore
}

func NewLoggingService(store persistence.LoggingStore) *LoggingService {
	return &LoggingService{
		store: store,
	}
}

func (s LoggingService) LoggInfo(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, INFO, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx, newLog)
}

func (s LoggingService) LoggError(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, ERROR, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx, newLog)
}

func (s LoggingService) LoggWarning(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, WARNING, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx, newLog)
}

func (s LoggingService) LoggSuccess(ctx context.Context, request *logging_service.LogRequest) (*logging_service.LogResult, error) {
	newLog := domain.NewLog(request.ServiceName, request.ServiceFunctionName, SUCCESS, request.UserID, request.IpAddress, request.Description)
	return s.store.SaveLog(ctx, newLog)
}
