package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
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

	v := validator.New()
	handler.ValidatePassword(ctx, v)
	handler.ValidateUsername(ctx, v)
	errV := v.Struct(user)
	if errV != nil {
		return &pb.RegisterResponse{
			Status: http.StatusNotAcceptable,
			UserID: "",
		}, errV
	}

	v := validator.New()
	handler.ValidatePassword(ctx, v)
	handler.ValidateUsername(ctx, v)
	errV := v.Struct(user)
	if errV != nil {
		return &pb.RegisterResponse{
			Status: http.StatusNotAcceptable,
			UserID: "",
		}, errV
	}

	userID, err := handler.service.Create(ctx, &user) //userID
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnauthorized,
			UserID: "",
		}, err
	}
	return &pb.RegisterResponse{
		Status: http.StatusCreated,
		UserID: userID,
	}, nil

}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := handler.service.GetByUsername(ctx, req.Username)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}
	var userForValidation domain.User
	userForValidation.Username = req.Data.GetUsername()
	userForValidation.Password = req.Data.GetPassword()

	v := validator.New()
	handler.ValidatePassword(ctx, v)
	handler.ValidateUsername(ctx, v)
	errV := v.Struct(userForValidation)
	if errV != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotAcceptable,
			Error:  "Bad credentials",
		}, errV
	}

	match := req.Password == user.Password

	if !match {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	token, _ := handler.Jwt.GenerateToken(user)

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

func (handler *AuthHandler) ValidatePassword(ctx context.Context, v *validator.Validate) {
	_ = v.RegisterValidation("password_validation", func(fl validator.FieldLevel) bool {
		if len(fl.Field().String()) < 8 {
			fmt.Println("Password must contain 8 characters or more!")
			return false
		}
		result, _ := regexp.MatchString("(.*[a-z].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain lower case characters!")
		}
		result, _ = regexp.MatchString("(.*[A-Z].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain upper case characters!")
		}
		result, _ = regexp.MatchString("(.*[0-9].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain numbers!")
		}

		result, _ = regexp.MatchString("(.*[!@#$%^&*(){}\\[:;\\]<>,\\.?~_+\\-\\\\=|/].*)", fl.Field().String())
		if !result {
			fmt.Println("Password must contain numbers!")
		}
		return result
	})

}

func (handler *AuthHandler) ValidateUsername(ctx context.Context, v *validator.Validate) {

	_ = v.RegisterValidation("username_validation", func(fl validator.FieldLevel) bool {
		if len(fl.Field().String()) < 4 && len(fl.Field().String()) > 16 {
			fmt.Println("Your username must be between 4 and 16 characters long.")
			return false
		}
		result, _ := regexp.MatchString("(.*[!@#$%^&*(){}\\[:;\\]<>,\\.?~_+\\-\\\\=|/].*)", fl.Field().String())
		if result {
			fmt.Println("Username contains special characters that are not allowed!")
		}
		return !result
	})

}
