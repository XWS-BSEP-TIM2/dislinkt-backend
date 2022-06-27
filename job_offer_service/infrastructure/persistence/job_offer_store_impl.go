package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobOfferDbStore struct {
	driverJobOffer *neo4j.Driver
}

func (store *JobOfferDbStore) Init() {
	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		errClear := clearGraphDB(transaction)
		if errClear != nil {
			return nil, errClear
		}
		errInit := initGraphDB(transaction)
		return nil, errInit
	})

	if err != nil {
		fmt.Println("JobOffer Graph Database INIT - Failed", err.Error())
	} else {
		fmt.Println("JobOffer Graph Database INIT - Successfully")
	}
}

func NewJobOfferDbStore(driver *neo4j.Driver) JobOfferStore {
	return &JobOfferDbStore{
		driverJobOffer: driver,
	}
}

func (store *JobOfferDbStore) Delete(ctx context.Context, jobId string) (bool, error) {
	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	r, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		if checkIfJobOfferExist(jobId, transaction) {
			_, err := transaction.Run(
				"MATCH (this_job:JOB) WHERE this_job.Id=$jobId DETACH DELETE this_job ",
				map[string]interface{}{"jobId": jobId})
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return false, err
	} else {
		return r.(bool), nil
	}
}

func (store *JobOfferDbStore) GetUserJobOffers(ctx context.Context, id primitive.ObjectID) ([]*domain.JobOffer, error) {
	/*
		filter := bson.M{"user_id": id}
		result, err := store.filter(filter)
		if err != nil {
			return nil, err
		} else {
			return result, nil
		}
	*/
	return nil, nil
}

func (store *JobOfferDbStore) Update(ctx context.Context, jobOffer *domain.JobOffer) error {
	/*
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
	*/
	return nil
}

func (store *JobOfferDbStore) Get(ctx context.Context, id primitive.ObjectID) (*domain.JobOffer, error) {
	/*
		filter := bson.M{"_id": id}
		return store.filterOne(filter)

	*/
	return nil, nil
}

func (store *JobOfferDbStore) GetAll(ctx context.Context) ([]*domain.JobOffer, error) {
	/*
		filter := bson.D{{}}
		return store.filter(filter)
	*/

	return nil, nil
}

func (store *JobOfferDbStore) Insert(ctx context.Context, jobOffer *domain.JobOffer) error {
	/*
		_, err := store.jobOffers.InsertOne(context.TODO(), jobOffer)
		if err != nil {
			return err
		}

	*/
	return nil
}

func (store *JobOfferDbStore) filterOne(filter interface{}) (profile *domain.JobOffer, err error) {
	/*
		result := store.jobOffers.FindOne(context.TODO(), filter)
		err = result.Decode(&profile)
	*/
	return nil, nil
}

func (store *JobOfferDbStore) filter(filter interface{}) ([]*domain.JobOffer, error) {
	/*
		cursor, err := store.jobOffers.Find(context.TODO(), filter)
		defer cursor.Close(context.TODO())

		if err != nil {
			return nil, err
		}
		return decode(cursor)
	*/
	return nil, nil
}

func decode(cursor *mongo.Cursor) (profiles []*domain.JobOffer, err error) {
	/*
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
	*/
	return nil, nil
}

func (store *JobOfferDbStore) Search(ctx context.Context, search string) ([]*domain.JobOffer, error) {
	/*
		var jobOffers []*domain.JobOffer

		err := filter(store, search, "position", &jobOffers)
		if err != nil {
			return nil, err
		}

		return jobOffers, nil

	*/
	return nil, nil
}

func filter(store *JobOfferDbStore, searchPart string, paramName string, jobOffers *[]*domain.JobOffer) error {
	/*
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

	*/
	return nil
}

func appendJobOffer(destination *[]*domain.JobOffer, source *domain.JobOffer) {
	/*
		for _, user := range *destination {
			if user.Id == source.Id {
				return
			}
		}
		*destination = append(*destination, source)
	*/
}

func (store *JobOfferDbStore) DeleteAll(ctx context.Context) {
	/*
		store.jobOffers.DeleteMany(context.TODO(), bson.D{{}})
	*/
}
