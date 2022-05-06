package persistence

import (
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"time"
)

const (
	DATABASE   = "post_db"
	COLLECTION = "post"
)

type ConnectionDBStore struct {
	connectionDB *neo4j.Driver
}

func NewConnectionDBStore(client *neo4j.Driver) domain.ConnectionStore {
	return &ConnectionDBStore{
		connectionDB: client,
	}
}

func (store *ConnectionDBStore) Init() {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
		fmt.Println("Connection Graph Database INIT FAILED!!! ", err.Error())
	} else {
		fmt.Println("Connection Graph Database INIT")
	}

}

func (store *ConnectionDBStore) GetFriends(userID string) ([]domain.UserConn, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	friends, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (this_user:USER) -[:FRIEND]-> (my_friend:USER) WHERE this_user.userID=$uID RETURN my_friend.userID, my_friend.isPrivate",
			map[string]interface{}{"uID": userID})

		if err != nil {
			return nil, err
		}

		var friends []domain.UserConn
		for result.Next() {
			friends = append(friends, domain.UserConn{UserID: result.Record().Values[0].(string), IsPublic: result.Record().Values[1].(bool)})
		}
		return friends, nil

	})
	if err != nil {
		return nil, err
	}

	return friends.([]domain.UserConn), nil
}

func (store *ConnectionDBStore) GetBlockeds(userID string) ([]domain.UserConn, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	blockedUsers, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (this_user:USER) -[:BLOCK]-> (my_friend:USER) WHERE this_user.userID=$uID RETURN my_friend.userID, my_friend.isPrivate",
			map[string]interface{}{"uID": userID})

		if err != nil {
			return nil, err
		}

		var friends []domain.UserConn
		for result.Next() {
			friends = append(friends, domain.UserConn{UserID: result.Record().Values[0].(string), IsPublic: result.Record().Values[1].(bool)})
		}
		return friends, nil

	})
	if err != nil {
		return nil, err
	}

	return blockedUsers.([]domain.UserConn), nil

}

func (store *ConnectionDBStore) Register(userID string, isPublic bool) (*pb.ActionResult, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{}

		if checkIfUserExist(userID, transaction) {
			actionResult.Status = 406
			actionResult.Msg = "error user with ID:" + userID + " already exist"
			return actionResult, nil
		}

		_, err := transaction.Run(
			"CREATE (new_user:USER{userID:$userID, isPrivate:$isPrivate})",
			map[string]interface{}{"userID": userID, "isPrivate": isPublic}) //TODO: promeniti u proto da bude isPrivate a ne isPublic, u Neo4J je isPrivate

		if err != nil {
			actionResult.Msg = "error while creating new node with ID:" + userID
			actionResult.Status = 501
			return actionResult, err
		}

		actionResult.Msg = "successfully created new node with ID:" + userID
		actionResult.Status = 201

		return actionResult, err
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}
func (store *ConnectionDBStore) AddFriend(userIDa, userIDb string) (*pb.ActionResult, error) {
	/*
				Dodavanje novog prijatelja je moguce ako:
		         - userA i userB postoji
				 - userA nije prijatelj sa userB
				 - userA nije blokirao userB
			   	 - userA nije blokiran od strane userB
	*/

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "userIDa is same as userIDb", Status: 400}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg", Status: 0}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfFriendExist(userIDa, userIDb, transaction) || checkIfFriendExist(userIDb, userIDa, transaction) {
				actionResult.Msg = "users are already friends"
				actionResult.Status = 400 //bad request
				return actionResult, nil
			} else {
				if checkIfBlockExist(userIDa, userIDb, transaction) || checkIfBlockExist(userIDb, userIDa, transaction) {
					actionResult.Msg = "block already exist"
					actionResult.Status = 400 //bad request
					return actionResult, nil
				} else {
					// TODO: dodati naknadno provu da li je profil privatan ili public

					dateNow := time.Now().Local().Unix()
					result, err := transaction.Run(
						"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
							"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
							"CREATE (u1)-[r1:FRIEND {date: $dateNow, msgID: $msgID}]->(u2) "+
							"CREATE (u2)-[r2:FRIEND {date: $dateNow, msgID: $msgID}]->(u1) "+
							"RETURN r1.date, r2.date",
						map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb, "dateNow": dateNow, "msgID": "newMsgID"})

					if err != nil || result == nil {
						actionResult.Msg = "error while creating new friends IDa:" + userIDa + " and IDb:" + userIDb
						actionResult.Status = 501
						return actionResult, err
					}

				}
			}
		} else {
			actionResult.Msg = "user does not exist"
			actionResult.Status = 400 //bad request
			return actionResult, nil
		}

		actionResult.Msg = "successfully created new friends IDa:" + userIDa + " and IDb:" + userIDb
		actionResult.Status = 201

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}

func (store *ConnectionDBStore) AddBlockUser(userIDa, userIDb string) (*pb.ActionResult, error) {

	/*
			UserA moze da blokira UserB ako:
			 - UserA nije vec blokirao UserB
		     - UserB vec nije blokirao prvi UserA
		  	Uspesno blokiranje rezultuje raskidanjem FRIEND veza izmedju ova dva cvora :(
	*/

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "userIDa is same as userIDb", Status: 400}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg", Status: 0}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfBlockExist(userIDa, userIDb, transaction) || checkIfBlockExist(userIDb, userIDa, transaction) {
				actionResult.Msg = "block already exist"
				actionResult.Status = 400 //bad request
				return actionResult, nil
			} else {
				if checkIfFriendExist(userIDa, userIDb, transaction) {
					removeFriend(userIDa, userIDb, transaction)
				}
				if checkIfFriendExist(userIDb, userIDa, transaction) {
					removeFriend(userIDb, userIDa, transaction)
				}
				blockUser(userIDa, userIDb, transaction)

				actionResult.Msg = "UserIDA:" + userIDa + " successfully blocked UserIDB:" + userIDb + " no longer friends!"
				actionResult.Status = 200
				return actionResult, nil
			}

		} else {
			actionResult.Msg = "user does not exist"
			actionResult.Status = 400 //bad request
			return actionResult, nil
		}
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}

func (store *ConnectionDBStore) RemoveFriend(userIDa, userIDb string) (*pb.ActionResult, error) {

	/*
		UserA mora biti prijatelj sa UserB (ne sme biti blokiran)
		UserA izbacuje prijatelja UserB, cepaju se obe prijateljske veze
	*/

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "userIDa is same as userIDb", Status: 400}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg", Status: 0}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfFriendExist(userIDa, userIDb, transaction) || checkIfFriendExist(userIDb, userIDa, transaction) {

				removeFriend(userIDa, userIDb, transaction)
				removeFriend(userIDb, userIDa, transaction)

			} else {
				actionResult.Msg = "users are not friends"
				actionResult.Status = 400 //bad request
				return actionResult, nil
			}
		} else {
			actionResult.Msg = "user does not exist"
			actionResult.Status = 400 //bad request
			return actionResult, nil
		}

		actionResult.Msg = "successfully user IDa:" + userIDa + " removed user IDb:" + userIDb
		actionResult.Status = 200

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}

func (store *ConnectionDBStore) UnblockUser(userIDa, userIDb string) (*pb.ActionResult, error) {
	/*
		UserA moze da unblokira useraB samo ako ga je on blokirao
	*/

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "userIDa is same as userIDb", Status: 400}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg", Status: 0}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfBlockExist(userIDb, userIDa, transaction) {
				actionResult.Msg = "UserB:" + userIDb + " first block UserA:" + userIDa
				actionResult.Status = 400 //bad request
				return actionResult, nil
			} else {
				if checkIfBlockExist(userIDa, userIDb, transaction) {
					if unblockUser(userIDa, userIDb, transaction) {
						actionResult.Msg = "successfully user IDa:" + userIDa + " unblock user IDb:" + userIDb
						actionResult.Status = 200
						return actionResult, nil
					}
				} else {
					actionResult.Msg = "UserA:" + userIDa + " and UserB:" + userIDb + " are nod blocked"
					actionResult.Status = 400 //bad request
					return actionResult, nil
				}
			}
		} else {
			actionResult.Msg = "user does not exist"
			actionResult.Status = 400 //bad request
			return actionResult, nil
		}

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}

func (store *ConnectionDBStore) GetRecommendation(userID string) ([]*domain.UserConn, error) {

	/*
		useru koji salje zahtevm, preporucicemo mu 20 prijatelje njegovih prijatelja
		ali necemo mu preporuciti one koje je on blokirao ili koji su njega blokirali

		takodje dobice jos do 20 preporuka ostlaih usera koji se ne nalaze u prvom skupu

		Metoda GetRecommendation vraca ukupno do 40 disjunktih preporuka
			- do 20 preporuka na osnovu prijatelja
			- do 20 preporuka random

	*/

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	recommendation, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		var recommendation []*domain.UserConn

		friendsOfFriends, err1 := getFriendsOfFriendsButNotBlockedRecommendation(userID, transaction)
		if err1 != nil {
			return recommendation, err1
		}

		for _, recommend := range friendsOfFriends {
			recommendation = append(recommendation, recommend)
		}

		famousRecom, err2 := getFriendRecommendation(userID, transaction)
		if err2 != nil {
			return recommendation, err2
		}

		var addNewRecommend bool = true
		for _, recommend := range famousRecom {
			addNewRecommend = true
			for _, r := range recommendation {
				if recommend.UserID == r.UserID {
					addNewRecommend = false
					break
				}
			}
			if addNewRecommend {
				recommendation = append(recommendation, recommend)
			}
		}

		return recommendation, err1

	})
	if err != nil || recommendation == nil {
		return nil, err
	}

	return recommendation.([]*domain.UserConn), nil
}
