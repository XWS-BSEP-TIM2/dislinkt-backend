package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "profile_db"
	COLLECTION = "profile"
)

type ProfileMongoDbStore struct {
	profiles *mongo.Collection
}

func (store *ProfileMongoDbStore) Update(profile *domain.Profile) error {

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

	_, err := store.profiles.UpdateOne(context.TODO(), profileToUpdate, updatedProfile)

	if err != nil {
		return err
	}
	return nil

}

func (store *ProfileMongoDbStore) Get(id primitive.ObjectID) (*domain.Profile, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *ProfileMongoDbStore) GetAll() ([]*domain.Profile, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *ProfileMongoDbStore) Insert(profile *domain.Profile) error {
	_, err := store.profiles.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}
	return nil
}

func NewProfileMongoDbStore(client *mongo.Client) ProfileStore {
	profiles := client.Database(DATABASE).Collection(COLLECTION)
	return &ProfileMongoDbStore{
		profiles: profiles,
	}
}

func (store *ProfileMongoDbStore) filterOne(filter interface{}) (profile *domain.Profile, err error) {
	result := store.profiles.FindOne(context.TODO(), filter)
	err = result.Decode(&profile)
	return
}

func (store *ProfileMongoDbStore) filter(filter interface{}) ([]*domain.Profile, error) {
	cursor, err := store.profiles.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (profiles []*domain.Profile, err error) {
	for cursor.Next(context.TODO()) {
		var profile domain.Profile
		err = cursor.Decode(&profile)
		if err != nil {
			return
		}
		profiles = append(profiles, &profile)
	}
	err = cursor.Err()
	return
}
