package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/startup/config"
	"github.com/google/uuid"
	myLogger "gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type LoggingDbStore struct {
	Con  *config.Config
	Logg *myLogger.Logger
}

func NewLoggingDbStore(c *config.Config, logg *myLogger.Logger) LoggingStore {
	return &LoggingDbStore{
		Con:  c,
		Logg: logg,
	}
}

func (l *LoggingDbStore) SaveLog(ctx context.Context, log *domain.Log) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "SaveLog")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	id := uuid.New()
	log.Id = id.String()
	fmt.Println(log.ToString())
	return l.SaveText(ctx2, log.ToString())
}

func (l *LoggingDbStore) Save(ctx context.Context, log string) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "Save")
	defer span.Finish()

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

func (l *LoggingDbStore) SaveText(ctx context.Context, log string) (*logging_service.LogResult, error) {
	span := tracer.StartSpanFromContext(ctx, "SaveText")
	defer span.Finish()

	result := &logging_service.LogResult{Msg: "Error", Status: 0}
	_, err := l.Logg.Write([]byte(log + "\n"))
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	result.Msg = "Successfully created new Log"
	result.Status = 201
	return result, nil
}
