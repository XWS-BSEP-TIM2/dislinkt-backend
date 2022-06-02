package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobOfferService struct {
	store persistence.JobOfferStore
}

func NewJobOfferService(store persistence.JobOfferStore) *JobOfferService {
	return &JobOfferService{
		store: store,
	}
}

func (service *JobOfferService) Get(ctx context.Context, id primitive.ObjectID) (*domain.JobOffer, error) {
	return service.store.Get(ctx, id)
}

func (service *JobOfferService) GetAll(ctx context.Context) ([]*domain.JobOffer, error) {
	return service.store.GetAll(ctx)
}

func (service *JobOfferService) Insert(ctx context.Context, profile *domain.JobOffer) {
	service.store.Insert(ctx, profile)
}

func (service *JobOfferService) Update(ctx context.Context, jobOffer *domain.JobOffer) {
	service.store.Update(ctx, jobOffer)
}

func (service *JobOfferService) Search(ctx context.Context, search string) ([]*domain.JobOffer, error) {
	return service.store.Search(ctx, search)
}
