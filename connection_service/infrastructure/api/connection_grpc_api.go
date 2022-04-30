package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/application"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.ConnectionService
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		service: service,
	}
}

func (handler *ConnectionHandler) GetFriends(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {

	fmt.Println("[ConnectionHandler]:GetFriends")

	id := request.UserID
	friends, err := handler.service.GetFriends(id)
	if err != nil {
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range friends {
		fmt.Println("User", id, "is friend with", user.UserID)
		response.Users = append(response.Users, mapUserConn(user))
	}
	return response, nil
}

func (handler *ConnectionHandler) GetBlockeds(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	fmt.Println("[ConnectionHandler]:GetBlockeds")

	id := request.UserID
	blockedUsers, err := handler.service.GetBlockeds(id)
	if err != nil {
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range blockedUsers {
		fmt.Println("User", id, "block", user.UserID)
		response.Users = append(response.Users, mapUserConn(user))
	}
	return response, nil
}
func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:Register")
	userID := request.User.UserID
	isPublic := request.User.IsPublic
	return handler.service.Register(userID, isPublic)
}
func (handler *ConnectionHandler) AddFriend(ctx context.Context, request *pb.AddFriendRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:AddFriend")
	userIDa := request.AddFriendDTO.UserIDa
	userIDb := request.AddFriendDTO.UserIDb
	return handler.service.AddFriend(userIDa, userIDb)
}
func (handler *ConnectionHandler) AddBlockUser(ctx context.Context, request *pb.AddBlockUserRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:AddBlockUser")
	userIDa := request.AddBlockUserDTO.UserIDa
	userIDb := request.AddBlockUserDTO.UserIDb
	return handler.service.AddBlockUser(userIDa, userIDb)
}
