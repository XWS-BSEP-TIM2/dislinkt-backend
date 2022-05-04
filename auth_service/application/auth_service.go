package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/api/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/services"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct {
	store                 domain.UserStore
	profileServiceAddress string
}

func NewAuthService(store domain.UserStore, profileServiceAddress string) *AuthService {
	return &AuthService{
		store:                 store,
		profileServiceAddress: profileServiceAddress,
	}
}

func (service *AuthService) Create(ctx context.Context, user *domain.User, dto *dto.CreateProfileDto) error {
	err, id := service.store.Insert(ctx, user)
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
	_, err = profileClient.CreateProfile(ctx, &profileService.CreateProfileRequest{Profile: &profile})
	if err != nil {
		return err
	}
	return nil
}

func (service *AuthService) Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	return service.store.Get(ctx, id)
}

func (service *AuthService) GetAll(ctx context.Context) ([]*domain.User, error) {
	return service.store.GetAll(ctx)
}

func (service *AuthService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return service.store.GetByUsername(ctx, username)
}
