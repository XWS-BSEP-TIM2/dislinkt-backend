package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const (
	DATABASE   = "profile_db"
	COLLECTION = "profile"
)

type ProfileMongoDbStore struct {
	profiles *mongo.Collection
}

func (store *ProfileMongoDbStore) Update(ctx context.Context, profile *domain.Profile) error {

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

func (store *ProfileMongoDbStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.Profile, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *ProfileMongoDbStore) GetAll(ctx context.Context) ([]*domain.Profile, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *ProfileMongoDbStore) Insert(ctx context.Context, profile *domain.Profile) error {
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

func (store *ProfileMongoDbStore) Search(ctx context.Context, search string) ([]*domain.Profile, error) {
	var profiles []*domain.Profile
	search = strings.TrimSpace(search)
	splitedSearch := strings.Split(search, " ")
	for _, searchPart := range splitedSearch {
		err := filter(store, searchPart, "username", &profiles)
		if err != nil {
			return nil, err
		}
		err = filter(store, searchPart, "name", &profiles)
		if err != nil {
			return nil, err
		}
		err = filter(store, searchPart, "surname", &profiles)
		if err != nil {
			return nil, err
		}
	}

	return profiles, nil

}

func filter(store *ProfileMongoDbStore, searchPart string, paramName string, profiles *[]*domain.Profile) error {
	filteredProfiles, err := store.profiles.Find(context.TODO(), bson.M{paramName: primitive.Regex{Pattern: searchPart, Options: "i"}})
	if err != nil {
		return err
	}
	var filterResult []*domain.Profile
	err = filteredProfiles.All(context.TODO(), &filterResult)
	if err != nil {
		return err
	}
	for _, result := range filterResult {
		appendUser(profiles, result)
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

func (store *ProfileMongoDbStore) DeleteAll(ctx context.Context) {
	store.profiles.DeleteMany(context.TODO(), bson.D{{}})
}
