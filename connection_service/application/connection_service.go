package application

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/domain"
)

type ConnectionService struct {
	store domain.ConnectionStore
}

func NewConnectionService(store domain.ConnectionStore) *ConnectionService {
	return &ConnectionService{
		store: store,
	}
}

func (service *ConnectionService) GetFriends(id string) ([]*domain.UserConn, error) {

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetFriends(id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPublic: s.IsPublic})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) GetBlockeds(id string) ([]*domain.UserConn, error) {

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetBlockeds(id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPublic: s.IsPublic})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) Register(userID string, isPublic bool) (*pb.ActionResult, error) {
	return service.store.Register(userID, isPublic)
}

func (service *ConnectionService) AddFriend(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.AddFriend(userIDa, userIDb)
}

func (service *ConnectionService) AddBlockUser(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.AddBlockUser(userIDa, userIDb)
}

func (service *ConnectionService) RemoveFriend(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.RemoveFriend(userIDa, userIDb)
}

func (service *ConnectionService) UnblockUser(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.UnblockUser(userIDa, userIDb)
}

func (service *ConnectionService) GetRecommendation(userID string) ([]*domain.UserConn, error) {
	return service.store.GetRecommendation(userID)
}
