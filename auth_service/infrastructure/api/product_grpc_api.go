package api

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type UserHandler struct {
	pb.UnimplementedAuthServiceServer
	service *application.UserService
	Jwt     utils.JwtWrapper
}

func NewProductHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	var user domain.User
	user.Username = request.User.Username
	user.Password = request.User.Password
	println("Username", request.User.Username)
	handler.service.Create(&user)

	return &pb.RegisterResponse{
		Status: http.StatusCreated,
	}, nil

}

func (handler *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := handler.service.GetByUsername(req.UserData.Username)
	if err != nil {
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}
	//if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
	//	return &pb.LoginResponse{
	//		Status: http.StatusNotFound,
	//		Error:  "User not found",
	//	}, nil
	//}

	//match := utils.CheckPasswordHash(req.Password, user.Password)
	match := req.UserData.Password == user.Password

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
	}, nil
}

func (handler *UserHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := handler.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	user, err := handler.service.Get(getObjectId(claims.Id))
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}
	//if result := handler.service.Get.H.DB.Where(&models.User{Email: claims.Email}).First(&user); result.Error != nil {
	//	return &pb.ValidateResponse{
	//		Status: http.StatusNotFound,
	//		Error:  "User not found",
	//	}, nil
	//}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id.Hex(),
	}, nil
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
