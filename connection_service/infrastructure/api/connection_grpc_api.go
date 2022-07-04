package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "GetFriends")
	defer span.Finish()

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
	span := tracer.StartSpanFromContext(ctx, "GetBlockeds")
	defer span.Finish()

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

func (handler *ConnectionHandler) GetFriendRequests(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	span := tracer.StartSpanFromContext(ctx, "GetFriendRequests")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:GetFriendRequests")

	id := request.UserID
	friendRequests, err := handler.service.GetFriendRequests(id)
	if err != nil {
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range friendRequests {
		fmt.Println("User", id, "friend reqest to", user.UserID)
		response.Users = append(response.Users, mapUserConn(&user))
	}
	return response, nil
}

func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "Register")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:Register")
	userID := request.User.UserID
	isPrivate := request.User.IsPrivate
	return handler.service.Register(userID, isPrivate)
}
func (handler *ConnectionHandler) AddFriend(ctx context.Context, request *pb.AddFriendRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "AddFriend")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:AddFriend")
	userIDa := request.AddFriendDTO.UserIDa
	userIDb := request.AddFriendDTO.UserIDb
	return handler.service.AddFriend(userIDa, userIDb)
}
func (handler *ConnectionHandler) AddBlockUser(ctx context.Context, request *pb.AddBlockUserRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "AddBlockUser")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:AddBlockUser")
	userIDa := request.AddBlockUserDTO.UserIDa
	userIDb := request.AddBlockUserDTO.UserIDb
	return handler.service.AddBlockUser(userIDa, userIDb)
}

func (handler *ConnectionHandler) RemoveFriend(ctx context.Context, request *pb.RemoveFriendRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "RemoveFriend")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:RemoveFriend")
	userIDa := request.RemoveFriendDTO.UserIDa
	userIDb := request.RemoveFriendDTO.UserIDb
	return handler.service.RemoveFriend(userIDa, userIDb)
}

func (handler *ConnectionHandler) UnblockUser(ctx context.Context, request *pb.UnblockUserRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "UnblockUser")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:UnblockUser")
	userIDa := request.UnblockUserDTO.UserIDa
	userIDb := request.UnblockUserDTO.UserIDb
	return handler.service.UnblockUser(userIDa, userIDb)
}

func (handler *ConnectionHandler) GetRecommendation(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecommendation")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:GetRecommendation")

	id := request.UserID
	recommendation, err := handler.service.GetRecommendation(id)
	if err != nil {
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range recommendation {
		response.Users = append(response.Users, mapUserConn(user))
	}
	return response, nil
}

func (handler *ConnectionHandler) SendFriendRequest(ctx context.Context, request *pb.SendFriendRequestRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "SendFriendRequest")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:SendFriendRequest")
	userIDa := request.SendFriendRequestRequestDTO.UserIDa
	userIDb := request.SendFriendRequestRequestDTO.UserIDb
	return handler.service.SendFriendRequest(userIDa, userIDb)
}
func (handler *ConnectionHandler) UnsendFriendRequest(ctx context.Context, request *pb.UnsendFriendRequestRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "UnsendFriendRequest")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:UnsendFriendRequest")
	userIDa := request.UnsendFriendRequestRequestDTO.UserIDa
	userIDb := request.UnsendFriendRequestRequestDTO.UserIDb
	return handler.service.UnsendFriendRequest(userIDa, userIDb)
}

func (handler *ConnectionHandler) GetConnectionDetail(ctx context.Context, request *pb.GetConnectionDetailRequest) (*pb.ConnectionDetail, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnectionDetail")
	defer span.Finish()

	fmt.Println("[ConnectionHandler]:GetConnectionDetail")
	userIDa := request.UserIDa
	userIDb := request.UserIDb
	return handler.service.GetConnectionDetail(userIDa, userIDb)
}

func (handler *ConnectionHandler) ChangePrivacy(ctx context.Context, request *pb.ChangePrivacyRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "ChangePrivacy")
	defer span.Finish()

	userID := request.ChangePrivacyBody.UserID
	isPrivate := request.ChangePrivacyBody.IsPrivate
	return handler.service.ChangePrivacy(userID, isPrivate)
}

func (handler *ConnectionHandler) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.ContactsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetMyContacts")
	defer span.Finish()

	return handler.service.GetMyContacts(ctx, request)
}
