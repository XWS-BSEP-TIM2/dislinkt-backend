package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*domain.Profile, error)
	GetAll(ctx context.Context) ([]*domain.Profile, error)
	Insert(ctx context.Context, profile *domain.Profile) error
	Update(ctx context.Context, profile *domain.Profile) error
	Search(ctx context.Context, search string) ([]*domain.Profile, error)
}
