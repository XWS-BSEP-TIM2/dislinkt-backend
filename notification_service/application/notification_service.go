package application

import (
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application/adapters"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/startup/config"
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
