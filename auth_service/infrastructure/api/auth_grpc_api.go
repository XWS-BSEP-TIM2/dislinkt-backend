package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	service *application.AuthService
	Jwt     utils.JwtWrapper
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (handler *AuthHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	var user domain.User
	user1, _ := handler.service.GetByUsername(ctx, request.Username)
	if user1 != nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  "Username is not unique",
			UserID: "",
		}, nil
	}
	user.Username = request.GetUsername()
	user.Password = request.GetPassword()
	user.Email = request.GetEmail()

	token, err := utils.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	user.VerificationCode = token
	user.VerificationCodeTime = time.Now()

	userID, err := handler.service.Create(ctx, &user) //userID
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnauthorized,
			UserID: "",
		}, err
	}

	errSendVerification := handler.service.SendVerification(ctx, &user)
	if errSendVerification != nil {
		fmt.Println("Error:", errSendVerification.Error())
	}

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
		UserID: userID,
	}, nil

}

func (handler *AuthHandler) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	return handler.service.Verify(ctx, req.Username, req.Code)
}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := handler.service.GetByUsername(ctx, req.Username)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Username or password is incorrect",
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

	match := req.Password == user.Password
	if !match {
		user.NumOfErrTryLogin += 1
		user.LastErrTryLoginTime = time.Now()
		if user.NumOfErrTryLogin >= 6 {
			user.Locked = true
			user.LockReason = "your account is locked, due to many incorrect login attempts"
		}
		handler.service.Update(ctx, user)
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Username or password is incorrect",
		}, nil
	}

	token, _ := handler.Jwt.GenerateToken(user)

	user.NumOfErrTryLogin = 0

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

	user, err := handler.service.Get(ctx, getObjectId(claims.Id))
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
		Id:       claims.Id,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil

}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
