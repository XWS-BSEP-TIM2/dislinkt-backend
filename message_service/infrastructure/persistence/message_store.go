package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageStore interface {
	GetChat(ctx context.Context, id primitive.ObjectID) (*domain.Chat, error)
	Insert(ctx context.Context, profile *domain.Chat) error
	DeleteAll(ctx context.Context)
}
