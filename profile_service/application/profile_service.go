package application

import (
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

func (service *ProfileService) Get(id primitive.ObjectID) (*domain.Profile, error) {
	return service.store.Get(id)
}

func (service *ProfileService) GetAll() ([]*domain.Profile, error) {
	return service.store.GetAll()
}

func (service *ProfileService) Insert(profile *domain.Profile) {
	service.store.Insert(profile)

}

func (service *ProfileService) Update(profile *domain.Profile) {
	service.store.Update(profile)
}

func (service *ProfileService) Search(search string) ([]*domain.Profile, error) {
	return service.store.Search(search)
}
