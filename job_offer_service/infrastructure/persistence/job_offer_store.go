package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobOfferStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*domain.JobOffer, error)
	GetAll(ctx context.Context) ([]*domain.JobOffer, error)
	Insert(ctx context.Context, profile *domain.JobOffer) error
	Update(ctx context.Context, profile *domain.JobOffer) error
	Search(ctx context.Context, search string) ([]*domain.JobOffer, error)
	DeleteAll(ctx context.Context)
	GetUserJobOffers(ctx context.Context, id primitive.ObjectID) ([]*domain.JobOffer, error)
	Delete(ctx context.Context, id primitive.ObjectID) (int64, error)
}
