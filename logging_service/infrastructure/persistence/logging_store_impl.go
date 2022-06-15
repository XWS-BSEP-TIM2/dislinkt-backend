package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/startup/config"
	"github.com/google/uuid"
	"os"
)

type LoggingDbStore struct {
	Con *config.Config
}

func NewLoggingDbStore(c *config.Config) LoggingStore {
	return &LoggingDbStore{
		Con: c,
	}
}

func (l *LoggingDbStore) SaveLog(ctx context.Context, log *domain.Log) (*logging_service.LogResult, error) {
	id := uuid.New()
	log.Id = id.String()
	fmt.Println(log.ToString())
	return l.Save(ctx, log.ToString())
}

func (l *LoggingDbStore) Save(ctx context.Context, log string) (*logging_service.LogResult, error) {
	result := &logging_service.LogResult{Msg: "Error", Status: 0}

	f, err := os.OpenFile(l.Con.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(f)

	if _, err = f.WriteString(log + "\n"); err != nil {
		return result, err
	}

	result.Msg = "Successfully created new Log"
	result.Status = 201
	return result, nil
}
