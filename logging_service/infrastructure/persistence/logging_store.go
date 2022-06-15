package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
)

type LoggingStore interface {
	SaveLog(ctx context.Context, log *domain.Log) (*logging_service.LogResult, error)
}
