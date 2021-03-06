package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
)

type MessageStore interface {
	GetChat(ctx context.Context, msgID string) (*domain.Chat, error)
	Insert(ctx context.Context, profile *domain.Chat) (string, error)
	DeleteAll(ctx context.Context)
	Update(ctx context.Context, chat *domain.Chat) error
	UpdateWithMessages(ctx context.Context, chat *domain.Chat) error
}
