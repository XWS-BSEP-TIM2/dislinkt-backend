package application

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	notificationService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/startup/config"
)

type ConnectionService struct {
	store               domain.ConnectionStore
	NotificationService notificationService.NotificationServiceClient
	ProfileClient       profileService.ProfileServiceClient
}

func NewConnectionService(store domain.ConnectionStore, c *config.Config) *ConnectionService {
	return &ConnectionService{
		store:               store,
		NotificationService: NewNotificationClient(fmt.Sprintf("%s:%s", c.NotificationServiceHost, c.NotificationServicePort)),
		ProfileClient:       NewProfileClient(fmt.Sprintf("%s:%s", c.ProfileHost, c.ProfilePort)),
	}
}

func (service *ConnectionService) GetFriends(id string) ([]*domain.UserConn, error) {

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetFriends(id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPrivate: s.IsPrivate})
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
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPrivate: s.IsPrivate})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) GetFriendRequests(userID string) ([]domain.UserConn, error) {
	return service.store.GetFriendRequests(userID)
}

func (service *ConnectionService) Register(userID string, isPublic bool) (*pb.ActionResult, error) {
	return service.store.Register(userID, isPublic)
}

func (service *ConnectionService) DeleteUser(userID string) (*pb.ActionResult, error) {
	return service.store.DeleteUser(userID)
}

func (service *ConnectionService) AddFriend(userIDa, userIDb string) (*pb.ActionResult, error) {
	sender, _ := service.ProfileClient.Get(context.TODO(), &profileService.GetRequest{Id: userIDa})
	var notification notificationService.Notification
	notification.OwnerId = userIDb
	notification.ForwardUrl = "profile/" + userIDa
	notification.Text = "is now your friend"
	notification.UserFullName = sender.Profile.Name + " " + sender.Profile.Surname
	service.NotificationService.InsertNotification(context.TODO(), &notificationService.InsertNotificationRequest{Notification: &notification})
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

func (service *ConnectionService) SendFriendRequest(userIDa, userIDb string) (*pb.ActionResult, error) {
	sender, _ := service.ProfileClient.Get(context.TODO(), &profileService.GetRequest{Id: userIDa})
	var notification notificationService.Notification
	notification.OwnerId = userIDb
	notification.ForwardUrl = "profile/" + userIDb + "/requests"
	notification.Text = "sent you a friend request"
	notification.UserFullName = sender.Profile.Name + " " + sender.Profile.Surname
	service.NotificationService.InsertNotification(context.TODO(), &notificationService.InsertNotificationRequest{Notification: &notification})

	return service.store.SendFriendRequest(userIDa, userIDb)
}

func (service *ConnectionService) UnsendFriendRequest(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.UnsendFriendRequest(userIDa, userIDb)
}

func (service *ConnectionService) GetConnectionDetail(userIDa, userIDb string) (*pb.ConnectionDetail, error) {
	return service.store.GetConnectionDetail(userIDa, userIDb)
}

func (service *ConnectionService) ChangePrivacy(userID string, private bool) (*pb.ActionResult, error) {
	return service.store.ChangePrivacy(userID, private)
}

func (service *ConnectionService) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.ContactsResponse, error) {
	return service.store.GetMyContacts(ctx, request)
}
