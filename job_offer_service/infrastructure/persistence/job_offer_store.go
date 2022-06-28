package persistence

import (
	"context"
	joboffer_service "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
)

type JobOfferStore interface {
	Get(ctx context.Context, jobId string) (*domain.JobOffer, error)
	GetAll(ctx context.Context) ([]*domain.JobOffer, error)
	Insert(ctx context.Context, profile *domain.JobOffer) error
	Update(ctx context.Context, profile *domain.JobOffer) (bool, error)
	Search(ctx context.Context, search string) ([]*domain.JobOffer, error)
	GetUserJobOffers(ctx context.Context, userID string) ([]*domain.JobOffer, error)
	Delete(ctx context.Context, jobId string) (bool, error)
	Init()
	CreateUser(ctx context.Context, userID string) (*joboffer_service.ActionResult, error)
	UpdateUserSkills(ctx context.Context, userID string, skills []string) (*joboffer_service.ActionResult, error)
}
