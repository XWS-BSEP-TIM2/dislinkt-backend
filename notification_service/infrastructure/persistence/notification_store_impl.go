package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson"
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

func (store *NotificationMongoDbStore) MarkAllAsSeen(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (store *NotificationMongoDbStore) DeleteAll(ctx context.Context) {
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
