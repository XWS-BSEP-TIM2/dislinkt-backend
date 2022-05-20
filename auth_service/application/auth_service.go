package application

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AuthService struct {
	store                 domain.UserStore
	profileServiceAddress string
	emailService          *EmailService
}

func NewAuthService(store domain.UserStore, profileServiceAddress string, emailService *EmailService) *AuthService {
	return &AuthService{
		store:                 store,
		profileServiceAddress: profileServiceAddress,
		emailService:          emailService,
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

func (service *AuthService) SendVerification(ctx context.Context, user *domain.User) error {
	fmt.Println("DOSLI SMO U METODU AuthService:SendVerification", user.Email)
	return service.emailService.SendVerificationEmail(user.Email, user.Username, user.VerificationCode)
}

func (service *AuthService) Update(ctx context.Context, user *domain.User) error {
	return service.store.Update(ctx, user)
}

func (service *AuthService) Verify(ctx context.Context, username string, code string) (*authService.VerifyResponse, error) {
	user, err := service.store.GetByUsername(ctx, username)
	if err != nil {
		return &authService.VerifyResponse{Verified: false, Msg: "User not found"}, err
	}

	if user.Verified {
		return &authService.VerifyResponse{Verified: true, Msg: "The user has already been verified"}, nil
	}

	if user.VerificationCodeTime.Add(10 * time.Minute).Before(time.Now()) {
		return &authService.VerifyResponse{Verified: false, Msg: "The verification code is no longer valid"}, nil
	}

	if user.VerificationCode == code {
		user.Verified = true
		errUpdate := service.store.Update(ctx, user)
		if errUpdate != nil {
			return &authService.VerifyResponse{Verified: false, Msg: "error"}, errUpdate
		}
		return &authService.VerifyResponse{Verified: true, Msg: "you have successfully verified your account"}, nil
	}
	return &authService.VerifyResponse{Verified: false, Msg: "error"}, nil
}
