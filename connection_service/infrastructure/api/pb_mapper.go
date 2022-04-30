package api

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/domain"
)

func mapUserConn(userConn *domain.UserConn) *pb.User {
	userConnPb := &pb.User{
		UserID:   userConn.UserID,
		IsPublic: userConn.IsPublic,
	}

	return userConnPb
}
