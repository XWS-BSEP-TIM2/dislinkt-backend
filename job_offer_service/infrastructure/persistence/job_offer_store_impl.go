package persistence

import (
	"context"
	"fmt"
	joboffer_service "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
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

func (store *JobOfferDbStore) GetRecommendationJobOffer(ctx context.Context, userID string) ([]*domain.JobOffer, error) {
	return store.getMany(ctx, "recommendation", userID)
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
		} else if getManyParam == "recommendation" {
			jobOfferIds, err = getRecommendationJobOfferIds(param, transaction)
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

func (store *JobOfferDbStore) Update(ctx context.Context, jobOffer *domain.JobOffer) (bool, error) {
	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		if checkIfJobOfferExist(jobOffer.Id, transaction) {
			updateResult, errU := updateJobOfferData(jobOffer, transaction)
			if errU != nil || updateResult == false {
				return false, errU
			}

			isSkillUpdated, err := updateSkillsForJobOffer(jobOffer, transaction)
			if err != nil || isSkillUpdated == false {
				return false, err
			}

			return true, nil
		}
		return false, nil
	})

	if err != nil || result == nil {
		return false, err
	}

	return result.(bool), nil
}

func (store *JobOfferDbStore) CreateUser(ctx context.Context, userID string) (*joboffer_service.ActionResult, error) {
	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	actionResult := &joboffer_service.ActionResult{Status: 0, Msg: ""}

	actRes, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		actRes := &joboffer_service.ActionResult{Status: 0, Msg: "Error"}
		if checkIfUserExist(userID, transaction) {
			createNewUser(userID, transaction)
			actRes.Msg = "Successfully created new user " + userID
			actRes.Status = 200
		}
		return actRes, nil
	})

	if err != nil {
		actionResult.Msg = err.Error()
		actionResult.Status = 400
		return actionResult, err
	}

	if actRes != nil {
		actionResult = actRes.(*joboffer_service.ActionResult)
	}

	return actionResult, nil
}

func (store *JobOfferDbStore) UpdateUserSkills(ctx context.Context, userID string, skills []string) (*joboffer_service.ActionResult, error) {
	session := (*store.driverJobOffer).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	actionResult := &joboffer_service.ActionResult{Status: 0, Msg: "Error"}

	actRes, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		actRes := &joboffer_service.ActionResult{Status: 0, Msg: "Error"}
		if checkIfUserExist(userID, transaction) {
			isSkillUpdated, err := updateSkillsForUser(userID, skills, transaction)
			if err != nil || isSkillUpdated == false {
				return actRes, err
			}

			actRes.Status = 200
			actRes.Msg = "Successfully updated skills for user " + userID
			return actRes, nil
		}
		return actRes, nil
	})

	if err != nil {
		actionResult.Msg = err.Error()
		actionResult.Status = 400
		return actionResult, err
	}

	if actRes != nil {
		actionResult = actRes.(*joboffer_service.ActionResult)
	}

	return actionResult, nil
}
