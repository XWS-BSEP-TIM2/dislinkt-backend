package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "messages_db"
	COLLECTION = "messages"
)

type MessageMongoDbStore struct {
	messages *mongo.Collection
}

func NewMessageMongoDbStore(client *mongo.Client) MessageStore {
	messages := client.Database(DATABASE).Collection(COLLECTION)
	return &MessageMongoDbStore{
		messages: messages,
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

func (store *MessageMongoDbStore) filterOne(filter interface{}) (chat *domain.Chat, err error) {
	result := store.messages.FindOne(context.TODO(), filter)
	err = result.Decode(&chat)
	return
}

func (store *MessageMongoDbStore) GetChat(ctx context.Context, msgID string) (*domain.Chat, error) {
	span := tracer.StartSpanFromContext(ctx, "GetChat")
	defer span.Finish()

	id := getIdFromHex(msgID)
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *MessageMongoDbStore) Insert(ctx context.Context, chat *domain.Chat) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert")
	defer span.Finish()

	result, err := store.messages.InsertOne(context.TODO(), chat)
	if err != nil {
		return "", err
	}
	chat.Id = result.InsertedID.(primitive.ObjectID)
	return chat.Id.Hex(), nil
}

func (store *MessageMongoDbStore) DeleteAll(ctx context.Context) {
	span := tracer.StartSpanFromContext(ctx, "DeleteAll")
	defer span.Finish()

	store.messages.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *MessageMongoDbStore) UpdateWithMessages(ctx context.Context, chat *domain.Chat) error {
	span := tracer.StartSpanFromContext(ctx, "UpdateWithMessages")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	chatToUpdate := bson.M{"_id": chat.Id}
	updatedChat := bson.M{"$set": bson.M{
		"userASeenDate": chat.UserASeenDate,
		"userBSeenDate": chat.UserBSeenDate,
		"messages":      chat.Messages,
	}}

	_, err := store.messages.UpdateOne(ctx2, chatToUpdate, updatedChat)

	if err != nil {
		return err
	}
	return nil
}

func (store *MessageMongoDbStore) Update(ctx context.Context, chat *domain.Chat) error {
	span := tracer.StartSpanFromContext(ctx, "Update")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	chatToUpdate := bson.M{"_id": chat.Id}
	updatedChat := bson.M{"$set": bson.M{
		"userASeenDate": chat.UserASeenDate,
		"userBSeenDate": chat.UserBSeenDate,
	}}

	_, err := store.messages.UpdateOne(ctx2, chatToUpdate, updatedChat)

	if err != nil {
		return err
	}
	return nil
}
