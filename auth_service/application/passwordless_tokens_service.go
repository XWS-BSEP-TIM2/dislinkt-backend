package application

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordlessTokenService struct {
	store        domain.PasswordlessTokenStore
	emailService *EmailService
}

func NewPasswordlessTokenService(store domain.PasswordlessTokenStore, emailService *EmailService) *PasswordlessTokenService {
	return &PasswordlessTokenService{
		store:        store,
		emailService: emailService,
	}
}

func (service *PasswordlessTokenService) Create(ctx context.Context, token *domain.PasswordlessToken) (string, error) {
	err, id := service.store.Insert(ctx, token)
	if err != nil {
		return "", err
	}
	return id, err
}

func (service *PasswordlessTokenService) Get(ctx context.Context, id primitive.ObjectID) (*domain.PasswordlessToken, error) {
	return service.store.Get(ctx, id)
}

func (service *PasswordlessTokenService) Delete(ctx context.Context, id primitive.ObjectID) {
	service.store.DeleteById(ctx, id)
}

func (service *PasswordlessTokenService) GetByTokenCode(ctx context.Context, tokenCode string) (*domain.PasswordlessToken, error) {
	return service.store.GetByTokenCode(ctx, tokenCode)
}

func (service *PasswordlessTokenService) SendMagicLink(ctx context.Context, user *domain.User, tokenCode string) {
	service.emailService.SendMagicLink(user.Email, tokenCode, "https://localhost:4200/magic-link-login")
}
