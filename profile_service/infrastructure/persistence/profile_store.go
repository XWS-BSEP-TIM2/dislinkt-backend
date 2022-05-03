package persistence

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileStore interface {
	Get(id primitive.ObjectID) (*domain.Profile, error)
	GetAll() ([]*domain.Profile, error)
	Insert(profile *domain.Profile) error
	Update(profile *domain.Profile) error
	Search(search string) ([]*domain.Profile, error)
}
