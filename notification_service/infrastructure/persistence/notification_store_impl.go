package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE            = "notification_db"
	COLLECTION          = "notifications"
	SETTINGS_COLLECTION = "user_settings"
)

type NotificationMongoDbStore struct {
	notifications *mongo.Collection
	userSettings  *mongo.Collection
}

func (store *NotificationMongoDbStore) GetOrInitUserSetting(ctx context.Context, userId primitive.ObjectID) *domain.UserSettings {
	settingsFilter := bson.M{"ownerId": userId}
	settings, _ := store.filterOneSettings(settingsFilter)
	if settings == nil {
		newSettings := domain.UserSettings{
			OwnerId:                 userId,
			PostNotifications:       true,
			ConnectionNotifications: true,
			MessageNotifications:    true,
		}

		err := store.InsertSetting(ctx, &newSettings)
		if err != nil {
			return nil
		}
		return &newSettings
	} else {
		return settings
	}
}

func (store *NotificationMongoDbStore) DeleteAllSettings(ctx context.Context) {
	store.userSettings.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *NotificationMongoDbStore) ModifyOrInsertSetting(ctx context.Context, setting *domain.UserSettings) {
	settingOld := store.GetOrInitUserSetting(ctx, setting.OwnerId)
	settingToUpdate := bson.M{"_id": settingOld.Id}
	updatedSetting := bson.M{"$set": bson.M{
		"postNotifications":       setting.PostNotifications,
		"connectionNotifications": setting.ConnectionNotifications,
		"messageNotifications":    setting.MessageNotifications,
	}}
	store.userSettings.UpdateOne(context.TODO(), settingToUpdate, updatedSetting)
}

func (store *NotificationMongoDbStore) InsertSetting(ctx context.Context, setting *domain.UserSettings) error {
	_, err := store.userSettings.InsertOne(context.TODO(), setting)
	if err != nil {
		return err
	}
	return nil
}

func (store *NotificationMongoDbStore) MarkAsSeen(ctx context.Context, notificationId primitive.ObjectID) {
	notificationToUpdate := bson.M{"_id": notificationId}
	updatedNotification := bson.M{"$set": bson.M{
		"seen": true,
	}}
	store.notifications.UpdateOne(context.TODO(), notificationToUpdate, updatedNotification)
}

func (store *NotificationMongoDbStore) DeleteAllNotifications(ctx context.Context) {
	store.notifications.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *NotificationMongoDbStore) filter(filter interface{}) ([]*domain.Notification, error) {
	cursor, err := store.notifications.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *NotificationMongoDbStore) filterOne(filter interface{}) (notification *domain.Notification, err error) {
	result := store.notifications.FindOne(context.TODO(), filter)
	err = result.Decode(&notification)
	return
}

func (store *NotificationMongoDbStore) filterOneSettings(filter interface{}) (settings *domain.UserSettings, err error) {
	result := store.userSettings.FindOne(context.TODO(), filter)
	err = result.Decode(&settings)
	return
}

func decode(cursor *mongo.Cursor) (notifications []*domain.Notification, err error) {
	for cursor.Next(context.TODO()) {
		var notification domain.Notification
		err = cursor.Decode(&notification)
		if err != nil {
			return
		}
		notifications = append(notifications, &notification)
	}
	err = cursor.Err()
	return
}

func (store *NotificationMongoDbStore) GetAll(ctx context.Context) ([]*domain.Notification, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *NotificationMongoDbStore) Insert(ctx context.Context, notification *domain.Notification) error {
	_, err := store.notifications.InsertOne(context.TODO(), notification)
	if err != nil {
		return err
	}
	return nil
}

func NewNotificationMongoDbStore(client *mongo.Client) NotificationStore {
	notificationsDb := client.Database(DATABASE).Collection(COLLECTION)
	settingsDb := client.Database(DATABASE).Collection(SETTINGS_COLLECTION)
	return &NotificationMongoDbStore{
		notifications: notificationsDb,
		userSettings:  settingsDb,
	}
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func getIdFromHex(userID string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(userID)
	return id
}
