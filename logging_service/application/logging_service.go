package application

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/persistence"
)

type LoggingService struct {
	store persistence.LoggingStore
}

func NewLoggingService(store persistence.LoggingStore) *LoggingService {
	return &LoggingService{
		store: store,
	}
}
