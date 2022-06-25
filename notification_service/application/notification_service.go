package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application/adapters"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/startup/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type NotificationService struct {
	store            persistence.NotificationStore
	ConnectionClient connectionService.ConnectionServiceClient
	ProfileClient    profileService.ProfileServiceClient
	LoggingService   loggingS.LoggingServiceClient
}

func NewNotificationService(store persistence.NotificationStore, c *config.Config, loggingService loggingS.LoggingServiceClient) *NotificationService {
	return &NotificationService{
		store:            store,
		ConnectionClient: adapters.NewConnectionClient(fmt.Sprintf("%s:%s", c.ConnectionHost, c.ConnectionPort)),
		ProfileClient:    adapters.NewProfileClient(fmt.Sprintf("%s:%s", c.ProfileHost, c.ProfilePort)),
		LoggingService:   loggingService,
	}
}

func (service *NotificationService) GetAllNotifications(ctx context.Context, request *pb.GetAllNotificationsRequest) (*pb.GetAllNotificationsResponse, error) {
	userId := request.UserID
	var userNotifications []*pb.Notification
	allNotifications, err := service.store.GetAll(ctx)

	if err == nil {
		for _, notification := range allNotifications {
			if notification.OwnerId.Hex() == userId {
				var newNotification = pb.Notification{
					OwnerId:    notification.OwnerId.Hex(),
					ForwardUrl: notification.ForwardUrl,
					Text:       notification.Text,
					Date:       &timestamppb.Timestamp{Seconds: notification.Date.Unix()},
					Seen:       notification.Seen,
				}
				userNotifications = append(userNotifications, &newNotification)
			}
		}
	}

	return &pb.GetAllNotificationsResponse{
		Notifications: userNotifications,
	}, err
}

func (service *NotificationService) MarkAllAsSeen(ctx context.Context, request *pb.MarkAllAsSeenRequest) (*pb.MarkAllAsSeenResponse, error) {
	userId := request.UserID
	allNotifications, err := service.store.GetAll(ctx)

	if err == nil {
		for _, notification := range allNotifications {
			if notification.OwnerId.Hex() == userId {
				if !notification.Seen {
					service.store.MarkAsSeen(ctx, notification.Id)
				}
			}
		}
	}

	return &pb.MarkAllAsSeenResponse{UserID: userId}, err
}

func (service *NotificationService) InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest) (*pb.InsertNotificationRequestResponse, error) {
	ownerId, err := primitive.ObjectIDFromHex(request.Notification.OwnerId)

	notification := &domain.Notification{
		OwnerId:    ownerId,
		ForwardUrl: request.Notification.ForwardUrl,
		Text:       request.Notification.Text,
		Date:       time.Now(),
		Seen:       false,
	}

	service.store.Insert(ctx, notification)

	return &pb.InsertNotificationRequestResponse{
		Notification: request.Notification,
	}, err
}
