package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	userService         *application.AuthService
	Jwt                 utils.JwtWrapper
	passwordlessService *application.PasswordlessTokenService
	apiTokenService     *application.ApiTokenService
}

func NewAuthHandler(service *application.AuthService, passwordlessServices *application.PasswordlessTokenService, apiTokenService *application.ApiTokenService) *AuthHandler {
	return &AuthHandler{
		userService:         service,
		passwordlessService: passwordlessServices,
		apiTokenService:     apiTokenService,
	}
}

func (handler *AuthHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	var user domain.User
	user1, _ := handler.userService.GetByUsername(ctx, request.Username)
	if user1 != nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  "Description is not unique",
			UserID: "",
		}, nil
	}
	user.Username = request.GetUsername()
	user.Password = utils.HashPassword(request.GetPassword())
	fmt.Println(user.Password)
	user.Email = request.GetEmail()

	token, err := utils.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	user.VerificationCode = token
	user.VerificationCodeTime = time.Now()

	userID, err := handler.userService.Create(ctx, &user) //userID
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnauthorized,
			UserID: "",
		}, err
	}

	errSendVerification := handler.userService.SendVerification(ctx, &user)
	if errSendVerification != nil {
		fmt.Println("Error:", errSendVerification.Error())
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
		UserID: userID,
	}, nil

}

func (handler *AuthHandler) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	return handler.userService.Verify(ctx, req.Username, req.Code)
}

func (handler *AuthHandler) ResendVerify(ctx context.Context, req *pb.ResendVerifyRequest) (*pb.ResendVerifyResponse, error) {
	return handler.userService.ResendVerify(ctx, req.Username)
}

func (handler *AuthHandler) Recovery(ctx context.Context, req *pb.RecoveryRequest) (*pb.RecoveryResponse, error) {
	return handler.userService.Recovery(ctx, req.Username)
}

func (handler *AuthHandler) Recover(ctx context.Context, req *pb.RecoveryRequestLogin) (*pb.LoginResponse, error) {
	response, err := handler.userService.Recover(ctx, req)
	if err != nil {
		return response, err
	}
	if response.Error != "" {
		return response, nil
	}
	return handler.Login(ctx, &pb.LoginRequest{Username: req.Username, Password: req.NewPassword})
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return handler.userService.ChangePassword(ctx, req)
}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := handler.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Description or password is incorrect",
		}, nil
	}

	if user.Locked {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  user.LockReason,
		}, nil
	}

	if !user.Verified {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  "Your Acc is not verified",
		}, nil
	}

	if user.NumOfErrTryLogin == 5 && !user.LastErrTryLoginTime.Add(1*time.Hour).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(1*time.Hour).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 4 && !user.LastErrTryLoginTime.Add(15*time.Minute).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(15*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 3 && !user.LastErrTryLoginTime.Add(3*time.Minute).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(3*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		user.NumOfErrTryLogin += 1
		user.LastErrTryLoginTime = time.Now()
		if user.NumOfErrTryLogin >= 6 {
			user.Locked = true
			user.LockReason = "your account is locked, due to many incorrect login attempts"
		}
		handler.userService.Update(ctx, user)
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Description or password is incorrect",
		}, nil
	}

	token, _ := handler.Jwt.GenerateToken(user)

	user.NumOfErrTryLogin = 0
	handler.userService.Update(ctx, user)

	return &pb.LoginResponse{
		Status:   http.StatusOK,
		Token:    token,
		Role:     domain.ConvertRoleToString(user.Role),
		Username: user.Username,
		UserID:   user.Id.Hex(),
	}, nil
}

func (handler *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := handler.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	user, err := handler.userService.Get(ctx, getObjectId(claims.Id))
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id.Hex(),
	}, nil
}

func (handler *AuthHandler) ExtractDataFromToken(ctx context.Context, req *pb.ExtractDataFromTokenRequest) (*pb.ExtractDataFromTokenResponse, error) {
	claims, err := handler.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ExtractDataFromTokenResponse{
			Id:       "",
			Username: "",
			Role:     "",
		}, err
	}

	return &pb.ExtractDataFromTokenResponse{
		Id:        claims.Id,
		Username:  claims.Username,
		Role:      claims.Role,
		TokenType: claims.TokenType,
		ApiCode:   claims.ApiCode,
	}, nil

}

func (handler *AuthHandler) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.LoginResponse, error) {
	token, err := handler.passwordlessService.GetByTokenCode(ctx, request.TokenCode)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Token does not exist",
		}, nil
	}
	if token.CreationDate.Add(15 * time.Minute).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Token is expired",
		}, err
	}

	handler.passwordlessService.Delete(ctx, token.Id)

	user, err := handler.userService.Get(ctx, token.UserId)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Description or password is incorrect",
		}, nil
	}

	if user.Locked {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  user.LockReason,
		}, nil
	}

	if !user.Verified {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  "Your Acc is not verified",
		}, nil
	}

	if user.NumOfErrTryLogin == 5 && !user.LastErrTryLoginTime.Add(1*time.Hour).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(1*time.Hour).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 4 && !user.LastErrTryLoginTime.Add(15*time.Minute).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(15*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 3 && !user.LastErrTryLoginTime.Add(3*time.Minute).Before(time.Now()) {
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(3*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	}

	tokenJwt, _ := handler.Jwt.GenerateToken(user)

	return &pb.LoginResponse{
		Status:   http.StatusOK,
		Token:    tokenJwt,
		Role:     domain.ConvertRoleToString(user.Role),
		Username: user.Username,
		UserID:   user.Id.Hex(),
	}, nil
}

func (handler *AuthHandler) SendEmailForPasswordlessLogin(ctx context.Context, request *pb.EmailForPasswordlessLoginRequest) (*pb.SendEmailForPasswordLoginResponse, error) {
	user, err := handler.userService.GetByEmail(ctx, request.Email)
	if err != nil || user == nil {
		return &pb.SendEmailForPasswordLoginResponse{
			Error: "Email does not exist",
		}, err
	}
	tokenCode, _ := utils.GenerateRandomStringURLSafe(30)
	token := domain.PasswordlessToken{
		TokenCode:    tokenCode,
		UserId:       user.Id,
		CreationDate: time.Now(),
	}
	handler.passwordlessService.Create(ctx, &token)
	handler.passwordlessService.SendMagicLink(ctx, user, tokenCode)
	return &pb.SendEmailForPasswordLoginResponse{
		Error: "",
	}, err
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func (handler *AuthHandler) GenerateApiToken(ctx context.Context, request *pb.ApiTokenRequest) (*pb.ApiTokenResponse, error) {
	user, _ := handler.userService.Get(ctx, getObjectId(request.UserId))
	if user == nil {
		error := pb.ErrorResponse{ErrorCode: 500, Message: "User ID does not exist"}
		return &pb.ApiTokenResponse{Error: &error}, nil
	}
	tokenCode, err := handler.apiTokenService.Create(ctx, getObjectId(request.UserId))
	if err != nil {
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Unable to create api token"}
		return &pb.ApiTokenResponse{Error: &error}, nil
	}
	return &pb.ApiTokenResponse{TokenCode: tokenCode, Error: nil}, nil
}

func (handler *AuthHandler) ValidateApiToken(ctx context.Context, request *pb.ValidateApiTokenRequest) (*pb.ValidateApiTokenResponse, error) {
	token, err := handler.apiTokenService.GetByTokenCode(ctx, request.TokenCode)
	if err != nil {
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Error while searching token"}
		return &pb.ValidateApiTokenResponse{Error: &error}, nil
	}
	if token == nil {
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Invalid token"}
		return &pb.ValidateApiTokenResponse{Error: &error}, nil
	}
	return &pb.ValidateApiTokenResponse{Error: nil}, nil
}
