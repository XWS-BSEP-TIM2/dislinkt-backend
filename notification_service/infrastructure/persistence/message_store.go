package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
)

type NotificationStore interface {
	GetAll(ctx context.Context) (*domain.Notification, error)
	Insert(ctx context.Context, profile *domain.Notification) error
}
