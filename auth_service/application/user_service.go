package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/api/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/services"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	store                 domain.UserStore
	profileServiceAddress string
}

func NewUserService(store domain.UserStore, profileServiceAddress string) *UserService {
	return &UserService{
		store:                 store,
		profileServiceAddress: profileServiceAddress,
	}
}

func (service *UserService) Create(user *domain.User, dto *dto.CreateProfileDto) error {
	err, id := service.store.Insert(user)
	if err != nil {
		return err
	}
	profileClient := services.NewProfileClient(service.profileServiceAddress)
	profile := profileService.Profile{
		Id:          id,
		Name:        dto.FirstName,
		Username:    user.Username,
		Email:       dto.Email,
		IsPrivate:   dto.IsPrivate,
		Surname:     dto.LastName,
		BirthDate:   dto.BirthDate,
		Gender:      dto.Gender,
		Skills:      []*profileService.Skill{},
		Experiences: []*profileService.Experience{},
	}
	_, err = profileClient.CreateProfile(context.TODO(), &profileService.CreateProfileRequest{Profile: &profile})
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.User, error) {
	return service.store.Get(id)
}

func (service *UserService) GetAll() ([]*domain.User, error) {
	return service.store.GetAll()
}

func (service *UserService) GetByUsername(username string) (*domain.User, error) {
	return service.store.GetByUsername(username)
}
