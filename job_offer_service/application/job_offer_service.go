package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	joboffer_service "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.Get(ctx2, jobId)
}

func (service *JobOfferService) GetAll(ctx context.Context) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAll(ctx2)
}

func (service *JobOfferService) Insert(ctx context.Context, jobOffer *domain.JobOffer) {
	span := tracer.StartSpanFromContext(ctx, "Insert")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	uniqueCode, err := utils.GenerateRandomString(24)
	if err != nil {
		return
	}
	jobOffer.JobOfferUniqueCode = uniqueCode
	jobOffer.Id = uniqueCode
	service.store.Insert(ctx2, jobOffer)
}

func (service *JobOfferService) Update(ctx context.Context, jobOffer *domain.JobOffer) {
	span := tracer.StartSpanFromContext(ctx, "Update")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	service.store.Update(ctx2, jobOffer)
}

func (service *JobOfferService) Search(ctx context.Context, search string) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "Search")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.Search(ctx2, search)
}

func (service *JobOfferService) GetUserJobOffers(ctx context.Context, userID string) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUserJobOffers")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetUserJobOffers(ctx2, userID)
}

func (service *JobOfferService) Delete(ctx context.Context, jobID string) (bool, error) {
	span := tracer.StartSpanFromContext(ctx, "Delete")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.Delete(ctx2, jobID)
}

func (service *JobOfferService) CreateUser(ctx context.Context, userID string) (*joboffer_service.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateUser")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.CreateUser(ctx2, userID)
}

func (service *JobOfferService) UpdateUserSkills(ctx context.Context, userID string, skills []string) (*joboffer_service.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSkills")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.UpdateUserSkills(ctx2, userID, skills)
}

func (service *JobOfferService) GetRecommendationJobOffer(ctx context.Context, userID string) ([]*domain.JobOffer, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecommendationJobOffer")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetRecommendationJobOffer(ctx2, userID)
}
