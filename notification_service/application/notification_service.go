package application

import (
	"context"
	"fmt"
	connectionService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	profileService "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application/adapters"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/startup/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/peer"
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
	span := tracer.StartSpanFromContext(ctx, "GetAllNotifications")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := request.UserID
	var userNotifications []*pb.Notification
	allNotifications, err := service.store.GetAll(ctx2)

	if err == nil {
		for _, notification := range allNotifications {
			if notification.OwnerId.Hex() == userId {
				var newNotification = pb.Notification{
					OwnerId:      notification.OwnerId.Hex(),
					ForwardUrl:   notification.ForwardUrl,
					Text:         notification.Text,
					Date:         &timestamppb.Timestamp{Seconds: notification.Date.Unix()},
					Seen:         notification.Seen,
					UserFullName: notification.UserFullName,
				}
				userNotifications = append(userNotifications, &newNotification)
			}
		}
	}

	service.logg(ctx2, "SUCCESS", "GetUserNotifications", request.UserID, "Getting all user's notifications.")
	return &pb.GetAllNotificationsResponse{
		Notifications: userNotifications,
	}, err
}

func (service *NotificationService) MarkAllAsSeen(ctx context.Context, request *pb.MarkAllAsSeenRequest) (*pb.MarkAllAsSeenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "MarkAllAsSeen")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := request.UserID
	allNotifications, err := service.store.GetAll(ctx2)

	if err == nil {
		for _, notification := range allNotifications {
			if notification.OwnerId.Hex() == userId {
				if !notification.Seen {
					service.logg(ctx2, "SUCCESS", "MarkNotificationsAsSeen", request.UserID, "Marking all unread notifications as seen.")
					service.store.MarkAsSeen(ctx2, notification.Id)
				}
			}
		}
	}

	return &pb.MarkAllAsSeenResponse{UserID: userId}, err
}

func (service *NotificationService) InsertNotification(ctx context.Context, request *pb.InsertNotificationRequest) (*pb.InsertNotificationRequestResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "InsertNotification")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	ownerId, err := primitive.ObjectIDFromHex(request.Notification.OwnerId)

	notification := &domain.Notification{
		OwnerId:      ownerId,
		ForwardUrl:   request.Notification.ForwardUrl,
		Text:         request.Notification.Text,
		Date:         time.Now(),
		Seen:         false,
		UserFullName: request.Notification.UserFullName,
	}

	if service.UserAcceptsNotification(notification) {
		service.logg(ctx2, "SUCCESS", "SendNotification", request.Notification.OwnerId, "Sending notification to user.")
		service.store.Insert(ctx2, notification)
	} else {
		service.logg(ctx2, "WARNING", "SendNotification", request.Notification.OwnerId, "Notification declined by user's settings.")
	}

	return &pb.InsertNotificationRequestResponse{
		Notification: request.Notification,
	}, err
}

func (service *NotificationService) GetUserSettings(ctx context.Context, request *pb.GetUserSettingsRequest) (*pb.GetUserSettingsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUserSettings")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := request.UserID
	id, err := primitive.ObjectIDFromHex(userId)
	if err == nil {
		settings := service.store.GetOrInitUserSetting(ctx, id)
		service.logg(ctx2, "SUCCESS", "GetUserSettings", request.UserID, "Fetching informations about user's notification settings.")
		return &pb.GetUserSettingsResponse{
			UserID:                  settings.OwnerId.Hex(),
			PostNotifications:       settings.PostNotifications,
			ConnectionNotifications: settings.ConnectionNotifications,
			MessageNotifications:    settings.MessageNotifications,
		}, nil
	} else {
		service.logg(ctx2, "ERROR", "GetUserSettings", request.UserID, "No such user.")
		return &pb.GetUserSettingsResponse{}, err
	}
}

func (service *NotificationService) UpdateUserSettings(ctx context.Context, request *pb.UpdateUserSettingsRequest) (*pb.GetUserSettingsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUserSettings")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := request.UserID
	id, err := primitive.ObjectIDFromHex(userId)
	if err == nil {
		settings := domain.UserSettings{
			OwnerId:                 id,
			PostNotifications:       true,
			ConnectionNotifications: true,
			MessageNotifications:    true,
		}

		if request.SettingsCode == "2" {
			settings.PostNotifications = false
			settings.ConnectionNotifications = false
			settings.MessageNotifications = false
		} else if request.SettingsCode == "3" {
			settings.PostNotifications = true
			settings.ConnectionNotifications = false
			settings.MessageNotifications = false
		} else if request.SettingsCode == "4" {
			settings.PostNotifications = false
			settings.ConnectionNotifications = true
			settings.MessageNotifications = false
		} else if request.SettingsCode == "5" {
			settings.PostNotifications = false
			settings.ConnectionNotifications = false
			settings.MessageNotifications = true
		} else if request.SettingsCode == "6" {
			settings.PostNotifications = true
			settings.ConnectionNotifications = false
			settings.MessageNotifications = true
		} else if request.SettingsCode == "7" {
			settings.PostNotifications = false
			settings.ConnectionNotifications = true
			settings.MessageNotifications = true
		} else if request.SettingsCode == "8" {
			settings.PostNotifications = true
			settings.ConnectionNotifications = true
			settings.MessageNotifications = false
		}

		service.store.ModifyOrInsertSetting(ctx2, &settings)
		return &pb.GetUserSettingsResponse{
			UserID: settings.OwnerId.Hex(),
		}, nil
	} else {
		return &pb.GetUserSettingsResponse{}, err
	}
}

func (service *NotificationService) UserAcceptsNotification(notification *domain.Notification) bool {
	settings := service.store.GetOrInitUserSetting(context.TODO(), notification.OwnerId)
	if notification.Text == "sent you a message" {
		return settings.MessageNotifications
	} else if notification.Text == "is now your friend" || notification.Text == "sent you a friend request" {
		return settings.ConnectionNotifications
	} else if notification.Text == "posted on their profile" || notification.Text == "commented on your post" {
		return settings.PostNotifications
	}

	return true
}

func (service *NotificationService) logg(ctx context.Context, logType, serviceFunctionName, userID, description string) {
	span := tracer.StartSpanFromContext(ctx, "logg")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	ipAddress := ""
	p, ok := peer.FromContext(ctx)
	if ok {
		ipAddress = p.Addr.String()
	}
	if logType == "ERROR" {
		service.LoggingService.LoggError(ctx2, &loggingS.LogRequest{ServiceName: "NOTIFICATION_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "SUCCESS" {
		service.LoggingService.LoggSuccess(ctx2, &loggingS.LogRequest{ServiceName: "NOTIFICATION_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "WARNING" {
		service.LoggingService.LoggWarning(ctx2, &loggingS.LogRequest{ServiceName: "NOTIFICATION_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "INFO" {
		service.LoggingService.LoggInfo(ctx2, &loggingS.LogRequest{ServiceName: "NOTIFICATION_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	}
}
