package persistence

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DATABASE   = "job_offer_db"
	COLLECTION = "job_offers"
)

type JobOfferMongoDbStore struct {
	jobOffers *mongo.Collection
}

func (store *JobOfferMongoDbStore) Delete(ctx context.Context, id primitive.ObjectID) (int64, error) {
	result, err := store.jobOffers.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (store *JobOfferMongoDbStore) GetUserJobOffers(ctx context.Context, id primitive.ObjectID) ([]*domain.JobOffer, error) {
	filter := bson.M{"user_id": id}
	result, err := store.filter(filter)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (store *JobOfferMongoDbStore) Update(ctx context.Context, jobOffer *domain.JobOffer) error {

	profileToUpdate := bson.M{"_id": jobOffer.Id}
	updatedProfile := bson.M{"$set": bson.M{
		"position":     jobOffer.Position,
		"seniority":    jobOffer.Seniority,
		"description":  jobOffer.Description,
		"company_name": jobOffer.CompanyName,
		"technologies": jobOffer.Technologies,
	}}

	_, err := store.jobOffers.UpdateOne(ctx, profileToUpdate, updatedProfile)

	if err != nil {
		return err
	}
	return nil

}

func (store *JobOfferMongoDbStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.JobOffer, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *JobOfferMongoDbStore) GetAll(ctx context.Context) ([]*domain.JobOffer, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *JobOfferMongoDbStore) Insert(ctx context.Context, jobOffer *domain.JobOffer) error {
	_, err := store.jobOffers.InsertOne(context.TODO(), jobOffer)
	if err != nil {
		return err
	}
	return nil
}

func NewJobOfferMongoDbStore(client *mongo.Client) JobOfferStore {
	profiles := client.Database(DATABASE).Collection(COLLECTION)
	return &JobOfferMongoDbStore{
		jobOffers: profiles,
	}
}

func (store *JobOfferMongoDbStore) filterOne(filter interface{}) (profile *domain.JobOffer, err error) {
	result := store.jobOffers.FindOne(context.TODO(), filter)
	err = result.Decode(&profile)
	return
}

func (store *JobOfferMongoDbStore) filter(filter interface{}) ([]*domain.JobOffer, error) {
	cursor, err := store.jobOffers.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func decode(cursor *mongo.Cursor) (profiles []*domain.JobOffer, err error) {
	for cursor.Next(context.TODO()) {
		var profile domain.JobOffer
		err = cursor.Decode(&profile)
		if err != nil {
			return
		}
		profiles = append(profiles, &profile)
	}
	err = cursor.Err()
	return
}

func (store *JobOfferMongoDbStore) Search(ctx context.Context, search string) ([]*domain.JobOffer, error) {
	var jobOffers []*domain.JobOffer

	err := filter(store, search, "position", &jobOffers)
	if err != nil {
		return nil, err
	}

	return jobOffers, nil

}

func filter(store *JobOfferMongoDbStore, searchPart string, paramName string, jobOffers *[]*domain.JobOffer) error {
	filteredOffers, err := store.jobOffers.Find(context.TODO(), bson.M{paramName: primitive.Regex{Pattern: searchPart, Options: "i"}})
	if err != nil {
		return err
	}
	var filterResult []*domain.JobOffer
	err = filteredOffers.All(context.TODO(), &filterResult)
	if err != nil {
		return err
	}
	for _, result := range filterResult {
		appendJobOffer(jobOffers, result)
	}
	return nil
}

func appendJobOffer(destination *[]*domain.JobOffer, source *domain.JobOffer) {
	for _, user := range *destination {
		if user.Id == source.Id {
			return
		}
	}
	*destination = append(*destination, source)
}

func (store *JobOfferMongoDbStore) DeleteAll(ctx context.Context) {
	store.jobOffers.DeleteMany(context.TODO(), bson.D{{}})
}