package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/infrastructure/persistence"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileService struct {
	store persistence.ProfileStore
}

func NewProfileService(store persistence.ProfileStore) *ProfileService {
	return &ProfileService{
		store: store,
	}
}

func (service *ProfileService) Get(ctx context.Context, id primitive.ObjectID) (*domain.Profile, error) {
	return service.store.Get(ctx, id)
}

func (service *ProfileService) GetAll(ctx context.Context) ([]*domain.Profile, error) {
	return service.store.GetAll(ctx)
}

func (service *ProfileService) Insert(ctx context.Context, profile *domain.Profile) {
	service.store.Insert(ctx, profile)

}

func (service *ProfileService) Update(ctx context.Context, profile *domain.Profile) {
	//TODO if username changes update it in auth service
	service.store.Update(ctx, profile)
}

func (service *ProfileService) Search(ctx context.Context, search string) ([]*domain.Profile, error) {
	return service.store.Search(ctx, search)
}