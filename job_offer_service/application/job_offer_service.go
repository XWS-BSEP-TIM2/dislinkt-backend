package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/infrastructure/persistence"
)

type JobOfferService struct {
	store persistence.JobOfferStore
}

func NewJobOfferService(store persistence.JobOfferStore) *JobOfferService {
	return &JobOfferService{
		store: store,
	}
}

func (service *JobOfferService) Get(ctx context.Context, jobId string) (*domain.JobOffer, error) {
	return service.store.Get(ctx, jobId)
}

func (service *JobOfferService) GetAll(ctx context.Context) ([]*domain.JobOffer, error) {
	return service.store.GetAll(ctx)
}

func (service *JobOfferService) Insert(ctx context.Context, jobOffer *domain.JobOffer) {
	uniqueCode, err := utils.GenerateRandomString(24)
	if err != nil {
		return
	}
	jobOffer.JobOfferUniqueCode = uniqueCode
	jobOffer.Id = uniqueCode
	service.store.Insert(ctx, jobOffer)
}

func (service *JobOfferService) Update(ctx context.Context, jobOffer *domain.JobOffer) {
	service.store.Update(ctx, jobOffer)
}

func (service *JobOfferService) Search(ctx context.Context, search string) ([]*domain.JobOffer, error) {
	return service.store.Search(ctx, search)
}

func (service *JobOfferService) GetUserJobOffers(ctx context.Context, userID string) ([]*domain.JobOffer, error) {
	return service.store.GetUserJobOffers(ctx, userID)
}

func (service *JobOfferService) Delete(ctx context.Context, jobID string) (bool, error) {
	return service.store.Delete(ctx, jobID)
}
