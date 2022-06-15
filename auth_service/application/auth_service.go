package application

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	authService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbLogg "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	dgoogauth "github.com/dgryski/dgoogauth"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/peer"
	"net/http"
	"net/url"
	qr "rsc.io/qr"
	"time"
)

type AuthService struct {
	store                 domain.UserStore
	profileServiceAddress string
	emailService          *EmailService
	LoggingService        pbLogg.LoggingServiceClient
}

func NewAuthService(store domain.UserStore, profileServiceAddress string, emailService *EmailService, loggingService pbLogg.LoggingServiceClient) *AuthService {
	return &AuthService{
		store:                 store,
		profileServiceAddress: profileServiceAddress,
		emailService:          emailService,
		LoggingService:        loggingService,
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
	return service.emailService.SendVerificationEmail(user.Email, user.Username, user.VerificationCode)
}

func (service *AuthService) Update(ctx context.Context, user *domain.User) error {
	return service.store.Update(ctx, user)
}

func (service *AuthService) Verify(ctx context.Context, username string, code string) (*authService.VerifyResponse, error) {
	user, err := service.store.GetByUsername(ctx, username)
	if err != nil {
		service.logg(ctx, "ERROR", "Verify", "", "User not found")
		return &authService.VerifyResponse{Verified: false, Msg: "User not found"}, err
	}

	if user.Verified {
		service.logg(ctx, "ERROR", "Verify", user.Id.Hex(), "The user has already been verified")
		return &authService.VerifyResponse{Verified: true, Msg: "The user has already been verified"}, nil
	}

	if user.VerificationCodeTime.Add(10 * time.Minute).Before(time.Now()) {
		service.logg(ctx, "ERROR", "Verify", user.Id.Hex(), "The verification code is no longer valid")
		return &authService.VerifyResponse{Verified: false, Msg: "The verification code is no longer valid"}, nil
	}

	if user.VerificationCode == code {
		user.Verified = true
		errUpdate := service.store.Update(ctx, user)
		if errUpdate != nil {
			return &authService.VerifyResponse{Verified: false, Msg: "error"}, errUpdate
		}
		service.logg(ctx, "SUCCESS", "Verify", user.Id.Hex(), "you have successfully verified your account")
		return &authService.VerifyResponse{Verified: true, Msg: "you have successfully verified your account"}, nil
	}
	service.logg(ctx, "ERROR", "Verify", user.Id.Hex(), "verification code did not match")
	return &authService.VerifyResponse{Verified: false, Msg: "error"}, nil
}

func (service *AuthService) Recovery(ctx context.Context, username string) (*authService.RecoveryResponse, error) {
	user, err := service.store.GetByUsername(ctx, username)
	if err != nil {
		service.logg(ctx, "ERROR", "Recovery", "", "User not found")
		return &authService.RecoveryResponse{Status: 1, Msg: "User not found"}, err
	}

	if !user.Verified {
		service.logg(ctx, "ERROR", "Recovery", user.Id.Hex(), "Recovery error: Your Acc is not verified")
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
		service.logg(ctx, "ERROR", "Recovery", user.Id.Hex(), "Error sending email")
		return &authService.RecoveryResponse{Status: 2, Msg: "Error sending email"}, errSendEmail
	}

	errUpdate := service.Update(ctx, user)
	if errUpdate != nil {
		return &authService.RecoveryResponse{Status: 3, Msg: "Error"}, errUpdate
	}

	service.logg(ctx, "SUCCESS", "Recovery", user.Id.Hex(), "Successfully sent recovery code")
	return &authService.RecoveryResponse{Status: 4, Msg: "Check your email, we sent you recovery code"}, nil
}

func (service *AuthService) Recover(ctx context.Context, req *authService.RecoveryRequestLogin) (*authService.LoginResponse, error) {
	if req.NewPassword != req.ConfirmNewPassword {
		return &authService.LoginResponse{Status: http.StatusBadRequest, Error: "passwords do not match"}, nil
	}

	user, err := service.store.GetByUsername(ctx, req.Username)
	if err != nil {
		service.logg(ctx, "ERROR", "Recover", "", "User not found")
		return &authService.LoginResponse{Status: http.StatusBadRequest, Error: "User not found"}, err
	}

	if user.RecoveryPasswordCodeTime.Add(5 * time.Minute).Before(time.Now()) {
		service.logg(ctx, "ERROR", "Recover", user.Id.Hex(), "The recovery code is no longer valid")
		return &authService.LoginResponse{Status: http.StatusNotAcceptable, Error: "The recovery code is no longer valid"}, nil
	}

	if user.RecoveryPasswordCode == req.RecoveryCode {
		if user.Locked {
			user.Locked = false
			user.LockReason = ""
		}
		user.NumOfErrTryLogin = 0
		user.Password = utils.HashPassword(req.NewPassword)
		service.Update(ctx, user)
		service.logg(ctx, "SUCCESS", "Recover", user.Id.Hex(), "Successfully recovered Acc")
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
		service.logg(ctx, "ERROR", "ResendVerify", user.Id.Hex(), "Error sending email")
		return &authService.ResendVerifyResponse{Msg: "error sending email"}, errSendEmail
	}

	service.logg(ctx, "SUCCESS", "ResendVerify", user.Id.Hex(), "Successfully send new verify code")
	return &authService.ResendVerifyResponse{Msg: "Check your email, we sent you verification link"}, err
}

func (service *AuthService) ChangePassword(ctx context.Context, req *authService.ChangePasswordRequest) (*authService.ChangePasswordResponse, error) {
	if req.NewPassword != req.ConfirmNewPassword {
		return &authService.ChangePasswordResponse{Status: http.StatusBadRequest, Msg: "passwords do not match"}, nil
	}

	user, err := service.store.GetByUsername(ctx, req.Username)
	if err != nil {
		return &authService.ChangePasswordResponse{Status: http.StatusBadRequest, Msg: "User not found"}, err
	}

	match := utils.CheckPasswordHash(req.OldPassword, user.Password)
	if !match {
		return &authService.ChangePasswordResponse{
			Status: http.StatusNotFound,
			Msg:    "Username or password is incorrect",
		}, nil
	}

	user.Password = utils.HashPassword(req.NewPassword)
	service.Update(ctx, user)

	return &authService.ChangePasswordResponse{
		Status: http.StatusOK,
		Msg:    "you have successfully change your password",
	}, nil
}

func (service *AuthService) GenerateQR2FA(ctx context.Context, userId primitive.ObjectID) ([]byte, error) {
	user, err := service.store.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	secret := make([]byte, 10)
	_, err = rand.Read(secret)
	if err != nil {
		panic(err)
	}

	user.TFASecret = base32.StdEncoding.EncodeToString(secret)
	service.store.Update(ctx, user)

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}

	URL.Path += "/" + url.PathEscape("Dislinkt") + ":" + url.PathEscape(user.Username)

	params := url.Values{}
	params.Add("secret", user.TFASecret)
	params.Add("issuer", "Dislinkt")

	URL.RawQuery = params.Encode()
	fmt.Printf("URL is %s\n", URL.String())

	code, err := qr.Encode(URL.String(), qr.Q)

	if err != nil {
		return nil, err
	}
	return code.PNG(), nil
}

func (service *AuthService) Verify2fa(ctx context.Context, userId primitive.ObjectID, code string) error {
	user, err := service.store.Get(ctx, userId)

	if err != nil {
		return err
	}

	otpc := &dgoogauth.OTPConfig{
		Secret:      user.TFASecret,
		WindowSize:  3,
		HotpCounter: 0,
		// UTC:         true,
	}
	val, err := otpc.Authenticate(code)
	if err != nil {
		return err
	}
	if !val {
		return errors.New("Not recognize code")
	}

	return nil

}

func (service *AuthService) logg(ctx context.Context, logType, serviceFunctionName, userID, description string) {
	ipAddress := ""
	p, ok := peer.FromContext(ctx)
	if ok {
		ipAddress = p.Addr.String()
	}
	if logType == "ERROR" {
		service.LoggingService.LoggError(ctx, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "SUCCESS" {
		service.LoggingService.LoggSuccess(ctx, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "WARNING" {
		service.LoggingService.LoggWarning(ctx, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "INFO" {
		service.LoggingService.LoggInfo(ctx, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	}
}
