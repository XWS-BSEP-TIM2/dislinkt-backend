package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application/adapters"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/startup/config"
	"log"
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
	log.Fatal("ERROR: 1")
	panic("unimplemented get all")
}

func (service *NotificationService) MarkAllAsSeen(ctx context.Context, request *pb.MarkAllAsSeenRequest) (*pb.MarkAllAsSeenResponse, error) {
	log.Fatal("ERROR: 2")
	panic("unimplemented mark all as seen")
}

func (service *NotificationService) InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest) (*pb.InsertNotificationRequestResponse, error) {
	log.Fatal("ERROR: 3")
	panic("unimplemented insert notification")
}
