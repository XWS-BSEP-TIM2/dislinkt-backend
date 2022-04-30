package domain

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
)

type ConnectionStore interface {
	GetFriends(id string) ([]UserConn, error)
	GetBlockeds(userID string) ([]UserConn, error)
	Register(userID string, isPublic bool) (*pb.ActionResult, error)
	AddFriend(userIDa, userIDb string) (*pb.ActionResult, error)
	AddBlockUser(userIDa, userIDb string) (*pb.ActionResult, error)
}
