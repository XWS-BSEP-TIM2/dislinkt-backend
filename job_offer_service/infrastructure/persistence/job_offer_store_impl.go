package persistence

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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

func (store *JobOfferDbStore) Get(ctx context.Context, jobId string) (*domain.JobOffer, error) {
	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	r, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		jobOffer, errJobOffer := getJobOffer(jobId, transaction)
		if errJobOffer != nil {
			fmt.Println(errJobOffer.Error())
			return nil, errJobOffer
		}
		return jobOffer, nil
	})

	if err != nil {
		return nil, err
	} else {
		if r != nil {
			return r.(*domain.JobOffer), nil
		} else {
			return nil, nil
		}
	}
}

func (store *JobOfferDbStore) GetUserJobOffers(ctx context.Context, userID string) ([]*domain.JobOffer, error) {
	return store.getMany(ctx, "GetUserJobOffers", userID)
}

func (store *JobOfferDbStore) GetAll(ctx context.Context) ([]*domain.JobOffer, error) {
	return store.getMany(ctx, "GetAll", "")
}

func (store *JobOfferDbStore) Search(ctx context.Context, search string) ([]*domain.JobOffer, error) {
	return store.getMany(ctx, "Search", search)
}

func (store *JobOfferDbStore) getMany(ctx context.Context, getManyParam, param string) ([]*domain.JobOffer, error) {

	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	r, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		var err error
		var jobOfferIds []string
		if getManyParam == "GetAll" {
			jobOfferIds, err = getAllJobOffersIds(transaction)
		} else if getManyParam == "GetUserJobOffers" {
			jobOfferIds, err = getUserJobOffersIds(param, transaction)
		} else if getManyParam == "Search" {
			jobOfferIds, err = searchJobOffersIds(param, transaction)
		}
		if err != nil {
			return nil, err
		}

		var userJobOffers []*domain.JobOffer
		for _, id := range jobOfferIds {
			jobOffer, errJobOffer := getJobOffer(id, transaction)
			if errJobOffer != nil {
				fmt.Println(errJobOffer.Error())
				continue
			}
			if jobOffer != nil {
				userJobOffers = append(userJobOffers, jobOffer)
			}
		}
		return userJobOffers, nil
	})

	if err != nil {
		return nil, err
	} else {
		if r != nil {
			return r.([]*domain.JobOffer), nil
		} else {
			return nil, nil
		}
	}
}

func (store *JobOfferDbStore) Insert(ctx context.Context, jobOffer *domain.JobOffer) error {

	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		if !checkIfJobOfferExist(jobOffer.Id, transaction) {
			errCreateJob := createNewJobOffer(jobOffer, transaction)
			if errCreateJob != nil {
				return nil, errCreateJob
			}
			for _, skill := range jobOffer.Technologies {
				if !checkIfSkillExist(skill, transaction) {
					createNewSkill(skill, transaction)
				}
				if !checkIfSkillIsPresentInJobOffer(jobOffer.Id, skill, transaction) {
					addSkillToJobOffer(jobOffer.Id, skill, transaction)
				}
			}
		}

		return nil, nil
	})

	return err
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
