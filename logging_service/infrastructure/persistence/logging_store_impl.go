package persistence

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/startup/config"
)

const (
	ERROR   = "ERROR"
	INFO    = "INFO"
	WARNING = "WARNING"
	SUCCESS = "SUCCESS"
)

type LoggingDbStore struct {
	Con *config.Config
}

func NewLoggingDbStore(c *config.Config) LoggingStore {
	return &LoggingDbStore{
		Con: c,
	}
}
