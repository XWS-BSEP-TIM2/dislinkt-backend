package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "notification_db"
	COLLECTION = "notifications"
)

type NotificationMongoDbStore struct {
	notifications *mongo.Collection
}

func (n NotificationMongoDbStore) GetAll(ctx context.Context) (*domain.Notification, error) {
	//TODO implement me
	panic("implement me")
}

func (n NotificationMongoDbStore) Insert(ctx context.Context, profile *domain.Notification) error {
	//TODO implement me
	panic("implement me")
}

func NewNotificationMongoDbStore(client *mongo.Client) NotificationStore {
	notificationsDb := client.Database(DATABASE).Collection(COLLECTION)
	return &NotificationMongoDbStore{
		notifications: notificationsDb,
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
