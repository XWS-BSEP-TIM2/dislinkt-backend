package application

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

func (service *AuthService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return service.store.GetByEmail(ctx, email)
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

func (service *AuthService) Recovery(ctx context.Context, username string) (*authService.RecoveryResponse, error) {
	user, err := service.store.GetByUsername(ctx, username)
	if err != nil {
		return &authService.RecoveryResponse{Status: 1, Msg: "User not found"}, err
	}

	if !user.Verified {
		return &authService.RecoveryResponse{Status: 5, Msg: "Recovery error: Your Acc is not verified"}, nil
	}

	recoveryCode, err := utils.GenerateRandomString(8)
	if err != nil {
		return nil, err
	}

	user.RecoveryPasswordCode = recoveryCode
	user.RecoveryPasswordCodeTime = time.Now()

	errSendEmail := service.emailService.SendRecoveryEmail(user.Email, user.Username, recoveryCode)
	if errSendEmail != nil {
		return &authService.RecoveryResponse{Status: 2, Msg: "Error sending email"}, errSendEmail
	}

	errUpdate := service.Update(ctx, user)
	if errUpdate != nil {
		return &authService.RecoveryResponse{Status: 3, Msg: "Error"}, errUpdate
	}

	return &authService.RecoveryResponse{Status: 4, Msg: "Check your email, we sent you recovery code"}, nil
}

func (service *AuthService) Recover(ctx context.Context, req *authService.RecoveryRequestLogin) (*authService.LoginResponse, error) {
	//TODO: validirati password regexom dal je dovoljno dobar? ...
	if req.NewPassword != req.ConfirmNewPassword {
		return &authService.LoginResponse{Status: http.StatusBadRequest, Error: "passwords do not match"}, nil
	}

	user, err := service.store.GetByUsername(ctx, req.Username)
	if err != nil {
		return &authService.LoginResponse{Status: http.StatusBadRequest, Error: "User not found"}, err
	}

	if user.RecoveryPasswordCodeTime.Add(5 * time.Minute).Before(time.Now()) {
		return &authService.LoginResponse{Status: http.StatusNotAcceptable, Error: "The recovery code is no longer valid"}, nil
	}

	if user.RecoveryPasswordCode == req.RecoveryCode {
		if user.Locked {
			user.Locked = false
			user.LockReason = ""
		}
		user.NumOfErrTryLogin = 0
		user.Password = req.NewPassword
		service.Update(ctx, user)
		return &authService.LoginResponse{Status: http.StatusOK, Error: ""}, nil
	}
	return &authService.LoginResponse{Error: "Error"}, nil
}

func (service *AuthService) ResendVerify(ctx context.Context, username string) (*authService.ResendVerifyResponse, error) {
	user, err := service.store.GetByUsername(ctx, username)
	if err != nil {
		return &authService.ResendVerifyResponse{Msg: "User not found"}, err
	}

	if user.Verified {
		return &authService.ResendVerifyResponse{Msg: "The user has already been verified"}, nil
	}

	token, errRandom := utils.GenerateRandomStringURLSafe(32)
	if errRandom != nil {
		panic(errRandom)
	}
	user.VerificationCode = token
	user.VerificationCodeTime = time.Now()
	service.Update(ctx, user)

	errSendEmail := service.emailService.SendVerificationEmail(user.Email, user.Username, user.VerificationCode)
	if errSendEmail != nil {
		return &authService.ResendVerifyResponse{Msg: "error sending email"}, errSendEmail
	}

	return &authService.ResendVerifyResponse{Msg: "Check your email, we sent you verification link"}, err
}
