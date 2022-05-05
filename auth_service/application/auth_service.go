package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
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

func (service *AuthService) Create(ctx context.Context, user *domain.User) (string, error) {
	err, id := service.store.Insert(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
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
