package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
)

type EventsStore interface {
	GetAll(ctx context.Context) ([]*domain.Event, error)
	Insert(ctx context.Context, event *domain.Event) error
	DeleteAll(todo context.Context)
}
