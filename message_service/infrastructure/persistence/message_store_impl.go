package persistence

import (
	"context"
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
	id := getIdFromHex(msgID)
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *MessageMongoDbStore) Insert(ctx context.Context, chat *domain.Chat) error {
	_, err := store.messages.InsertOne(context.TODO(), chat)
	if err != nil {
		return err
	}
	return nil
}

func (store *MessageMongoDbStore) DeleteAll(ctx context.Context) {
	store.messages.DeleteMany(context.TODO(), bson.D{{}})
}

/*
func (store *MessageMongoDbStore) Update(ctx context.Context, profile *domain.Profile) error {

	profileToUpdate := bson.M{"_id": profile.Id}
	updatedProfile := bson.M{"$set": bson.M{
		"name":        profile.Name,
		"surname":     profile.Surname,
		"username":    profile.Username,
		"email":       profile.Email,
		"biography":   profile.Biography,
		"gender":      profile.Gender,
		"phoneNumber": profile.PhoneNumber,
		"birthDate":   profile.BirthDate,
		"isPrivate":   profile.IsPrivate,
		"skills":      profile.Skills,
		"experiences": profile.Experiences,
	}}

	_, err := store.messages.UpdateOne(context.TODO(), profileToUpdate, updatedProfile)

	if err != nil {
		return err
	}
	return nil

}

func (store *MessageMongoDbStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.Profile, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *MessageMongoDbStore) GetAll(ctx context.Context) ([]*domain.Profile, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *MessageMongoDbStore) Insert(ctx context.Context, profile *domain.Profile) error {
	_, err := store.messages.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}
	return nil
}

func NewProfileMongoDbStore(client *mongo.Client) MessageStore {
	messages := client.Database(DATABASE).Collection(COLLECTION)
	return &MessageMongoDbStore{
		messages: messages,
	}
}

func (store *MessageMongoDbStore) filterOne(filter interface{}) (profile *domain.Profile, err error) {
	result := store.messages.FindOne(context.TODO(), filter)
	err = result.Decode(&profile)
	return
}

func (store *MessageMongoDbStore) filter(filter interface{}) ([]*domain.Profile, error) {
	cursor, err := store.messages.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (messages []*domain.Profile, err error) {
	for cursor.Next(context.TODO()) {
		var profile domain.Profile
		err = cursor.Decode(&profile)
		if err != nil {
			return
		}
		messages = append(messages, &profile)
	}
	err = cursor.Err()
	return
}

func (store *MessageMongoDbStore) Search(ctx context.Context, search string) ([]*domain.Profile, error) {
	var messages []*domain.Profile
	search = strings.TrimSpace(search)
	splitedSearch := strings.Split(search, " ")
	for _, searchPart := range splitedSearch {
		err := filter(store, searchPart, "username", &messages)
		if err != nil {
			return nil, err
		}
		err = filter(store, searchPart, "name", &messages)
		if err != nil {
			return nil, err
		}
		err = filter(store, searchPart, "surname", &messages)
		if err != nil {
			return nil, err
		}
	}

	return messages, nil

}

func filter(store *MessageMongoDbStore, searchPart string, paramName string, messages *[]*domain.Profile) error {
	filteredProfiles, err := store.messages.Find(context.TODO(), bson.M{paramName: primitive.Regex{Pattern: searchPart, Options: "i"}})
	if err != nil {
		return err
	}
	var filterResult []*domain.Profile
	err = filteredProfiles.All(context.TODO(), &filterResult)
	if err != nil {
		return err
	}
	for _, result := range filterResult {
		appendUser(messages, result)
	}
	return nil
}

func appendUser(destination *[]*domain.Profile, source *domain.Profile) {
	for _, user := range *destination {
		if user.Id == source.Id {
			return
		}
	}
	*destination = append(*destination, source)
}

func (store *MessageMongoDbStore) DeleteAll(ctx context.Context) {
	store.messages.DeleteMany(context.TODO(), bson.D{{}})
}


*/
