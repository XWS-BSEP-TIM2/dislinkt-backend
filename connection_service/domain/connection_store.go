package domain

import (
	"context"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
)

type ConnectionStore interface {
	GetFriends(id string) ([]UserConn, error)
	GetBlockeds(userID string) ([]UserConn, error)
	Register(userID string, isPublic bool) (*pb.ActionResult, error)
	AddFriend(userIDa, userIDb string) (*pb.ActionResult, error)
	AddBlockUser(userIDa, userIDb string) (*pb.ActionResult, error)
	RemoveFriend(userIDa, userIDb string) (*pb.ActionResult, error)
	UnblockUser(userIDa, userIDb string) (*pb.ActionResult, error)
	GetRecommendation(userID string) ([]*UserConn, error)
	Init()
	SendFriendRequest(userIDa, userIDb string) (*pb.ActionResult, error)
	UnsendFriendRequest(userIDa, userIDb string) (*pb.ActionResult, error)
	GetConnectionDetail(userIDa, userIDb string) (*pb.ConnectionDetail, error)
	GetFriendRequests(userID string) ([]UserConn, error)
	ChangePrivacy(userID string, private bool) (*pb.ActionResult, error)
	GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.ContactsResponse, error)
}
