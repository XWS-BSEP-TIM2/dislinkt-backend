package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationStore interface {
	GetAll(ctx context.Context) ([]*domain.Notification, error)
	Insert(ctx context.Context, notification *domain.Notification) error
	DeleteAllNotifications(ctx context.Context)
	MarkAsSeen(ctx context.Context, notificationId primitive.ObjectID)

	DeleteAllSettings(ctx context.Context)
	ModifyOrInsertSetting(ctx context.Context, setting *domain.UserSettings)
	InsertSetting(ctx context.Context, setting *domain.UserSettings) error
	GetOrInitUserSetting(ctx context.Context, userId primitive.ObjectID) *domain.UserSettings
}
