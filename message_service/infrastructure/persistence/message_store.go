package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
)

type MessageStore interface {
	GetChat(ctx context.Context, msgID string) (*domain.Chat, error)
	Insert(ctx context.Context, profile *domain.Chat) error
	DeleteAll(ctx context.Context)
}
