package persistence

import (
	"context"
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
			friends = append(friends, domain.UserConn{UserID: result.Record().Values[0].(string), IsPrivate: result.Record().Values[1].(bool)})
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
			friends = append(friends, domain.UserConn{UserID: result.Record().Values[0].(string), IsPrivate: result.Record().Values[1].(bool)})
		}
		return friends, nil

	})
	if err != nil {
		return nil, err
	}

	return blockedUsers.([]domain.UserConn), nil

}

func (store *ConnectionDBStore) GetFriendRequests(userID string) ([]domain.UserConn, error) {
	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	friendsRequest, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (this_user:USER) <-[:REQUEST]- (user_requester:USER) WHERE this_user.userID=$uID RETURN user_requester.userID, user_requester.isPrivate",
			map[string]interface{}{"uID": userID})

		if err != nil {
			return nil, err
		}

		var friendsRequest []domain.UserConn
		for result.Next() {
			friendsRequest = append(friendsRequest, domain.UserConn{UserID: result.Record().Values[0].(string), IsPrivate: result.Record().Values[1].(bool)})
		}
		return friendsRequest, nil

	})
	if err != nil {
		return nil, err
	}

	return friendsRequest.([]domain.UserConn), nil
}

func (store *ConnectionDBStore) Register(userID string, isPrivate bool) (*pb.ActionResult, error) {

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
			map[string]interface{}{"userID": userID, "isPrivate": isPrivate}) //TODO: promeniti u proto da bude isPrivate a ne isPublic, u Neo4J je isPrivate

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

					isPrivate, err := isUserPrivate(userIDb, transaction)
					if err != nil {
						fmt.Println(err.Error())
						actionResult.Msg = err.Error()
						actionResult.Status = 400 //bad request
						return actionResult, nil
					}
					if isPrivate {
						// ako je profil privatan, onda uspeva samo ako je ovaj profil vec poslao zahtev pre
						if checkIfFriendRequestExist(userIDb, userIDa, transaction) {
							//ok postoji zahtev, mozemo spajati
							removeFriendRequest(userIDb, userIDa, transaction)
						} else {
							actionResult.Msg = "error UserIDb:" + userIDb + " are private profile and do not exist friend request"
							actionResult.Status = 400 //bad request
							return actionResult, nil
						}
					}

					// dodavanje prijatelja
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

					// brisanje zahteva za prijateljstvo za slucaj kada:
					// userA je privatan, i dobija request od userB
					// userA nije odgovorio na request, al je promenio svoj nalog iz private u public
					// u slucaju da mogu da se dodaju kao prijatelji onda treba obrisati ove veze
					if checkIfFriendRequestExist(userIDa, userIDb, transaction) {
						removeFriendRequest(userIDa, userIDb, transaction)
					}
					if checkIfFriendRequestExist(userIDb, userIDa, transaction) {
						removeFriendRequest(userIDb, userIDa, transaction)
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
				if checkIfFriendRequestExist(userIDa, userIDb, transaction) {
					removeFriendRequest(userIDa, userIDb, transaction)
				}
				if checkIfFriendRequestExist(userIDb, userIDa, transaction) {
					removeFriendRequest(userIDb, userIDa, transaction)
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

func (store *ConnectionDBStore) SendFriendRequest(userIDa, userIDb string) (*pb.ActionResult, error) {
	/*

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
					if checkIfFriendRequestExist(userIDa, userIDb, transaction) || checkIfFriendRequestExist(userIDb, userIDa, transaction) {
						actionResult.Msg = "FriendRequest already exist"
						actionResult.Status = 400 //bad request
						return actionResult, nil
					} else {

						isPrivate, err := isUserPrivate(userIDb, transaction)
						if err != nil {
							fmt.Println(err.Error())
							actionResult.Msg = err.Error()
							actionResult.Status = 400 //bad request
							return actionResult, nil
						}

						if isPrivate {
							dateNow := time.Now().Local().Unix()
							result, err := transaction.Run(
								"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
									"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
									"CREATE (u1)-[r1:REQUEST {date: $dateNow}]->(u2) "+
									"RETURN r1.date",
								map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb, "dateNow": dateNow})

							if err != nil || result == nil {
								actionResult.Msg = "error while creating new friends IDa:" + userIDa + " and IDb:" + userIDb
								actionResult.Status = 501
								return actionResult, err
							}
						} else {
							actionResult.Msg = "error userIDb:" + userIDb + " is not private user"
							actionResult.Status = 400 //bad request
							return actionResult, nil
						}
					}
				}
			}
		} else {
			actionResult.Msg = "user does not exist"
			actionResult.Status = 400 //bad request
			return actionResult, nil
		}

		actionResult.Msg = "User UserIDa:" + userIDa + " successfully send request to UserIDb:" + userIDb
		actionResult.Status = 201

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}

func (store *ConnectionDBStore) UnsendFriendRequest(userIDa, userIDb string) (*pb.ActionResult, error) {
	/*
		UserA moze da povuce zahtev za prijateljstvo samo ako je poslao
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
				actionResult.Msg = "Users UserIDa:" + userIDa + " and UserIDb:" + userIDb + " are blocked"
				actionResult.Status = 400 //bad request
				return actionResult, nil
			} else {
				if checkIfFriendRequestExist(userIDa, userIDb, transaction) {
					removeFriendRequest(userIDa, userIDb, transaction)
					actionResult.Msg = "Users UserIDa:" + userIDa + " unsend friend request to UserIDb:" + userIDb
					actionResult.Status = 200 //bad request
					return actionResult, nil
				} else {
					actionResult.Msg = "friend request do not exist"
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

func (store *ConnectionDBStore) GetConnectionDetail(userIDa, userIDb string) (*pb.ConnectionDetail, error) {

	/*

	 */
	if userIDa == userIDb {
		return &pb.ConnectionDetail{Error: "userIDa is same as userIDb"}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		connectionDetail := &pb.ConnectionDetail{UserIDa: userIDa, UserIDb: userIDb}

		// userIDa is not logged in or irrelevant
		// used for checking if userIDb account is private
		if userIDa == "000000000000000000000000" {
			if !checkIfUserExist(userIDb, transaction) {
				connectionDetail.Error = "user does not exist"
				return connectionDetail, nil
			}
			isPrivate, err := isUserPrivate(userIDb, transaction)
			if err != nil {
				connectionDetail.Error = err.Error()
				return connectionDetail, err
			}
			connectionDetail.IsPrivate = isPrivate
			connectionDetail.Relation = "NO_RELATION"
			return connectionDetail, nil
		}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {

			isPrivate, err := isUserPrivate(userIDb, transaction)
			if err != nil {
				connectionDetail.Error = err.Error()
				return connectionDetail, err
			}

			connectionDetail.IsPrivate = isPrivate

			if checkIfBlockExist(userIDa, userIDb, transaction) {
				connectionDetail.Relation = "A_BLOCK_B"
				return connectionDetail, nil
			}
			if checkIfBlockExist(userIDb, userIDa, transaction) {
				connectionDetail.Relation = "B_BLOCK_A"
				return connectionDetail, nil
			}

			if checkIfFriendExist(userIDa, userIDb, transaction) || checkIfFriendExist(userIDb, userIDa, transaction) {
				connectionDetail.Relation = "CONNECTED"
				return connectionDetail, nil
			}

			if checkIfFriendRequestExist(userIDa, userIDb, transaction) {
				connectionDetail.Relation = "PENDING"
				return connectionDetail, nil
			}
			if checkIfFriendRequestExist(userIDb, userIDa, transaction) {
				connectionDetail.Relation = "ACCEPT"
				return connectionDetail, nil
			}

			connectionDetail.Relation = "NO_RELATION"

		} else {
			connectionDetail.Error = "user does not exist"
			return connectionDetail, nil
		}

		return connectionDetail, nil
	})

	if result == nil {
		return &pb.ConnectionDetail{Error: "error"}, err
	} else {
		return result.(*pb.ConnectionDetail), err
	}
}

func (store *ConnectionDBStore) ChangePrivacy(userID string, private bool) (*pb.ActionResult, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg", Status: 0}

		if checkIfUserExist(userID, transaction) {
			isPrivate, err := isUserPrivate(userID, transaction)
			if err != nil {
				actionResult.Msg = err.Error()
				return actionResult, err
			}

			if isPrivate != private {
				ok, err := setUserPrivate(userID, private, transaction)
				if err != nil {
					actionResult.Msg = err.Error()
					actionResult.Status = 400
					return nil, err
				}
				if !ok {
					actionResult.Msg = "error updating privacy"
					return actionResult, nil
				} else {
					actionResult.Msg = "successfully changed privacy"
					actionResult.Status = 200
					return actionResult, nil
				}
			} else {
				actionResult.Msg = "same privacy"
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

func (store *ConnectionDBStore) GetMyContacts(ctx context.Context, request *pb.GetMyContactsRequest) (*pb.ContactsResponse, error) {

	userID := request.UserID

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	contacts, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (this_user:USER) -[r:FRIEND]-> (my_friend:USER) WHERE this_user.userID=$uID RETURN my_friend.userID, r.msgID ",
			map[string]interface{}{"uID": userID})

		if err != nil {
			return nil, err
		}

		var contacts []*pb.Contact
		for result.Next() {
			contacts = append(contacts, &pb.Contact{UserID: result.Record().Values[0].(string), MsgID: result.Record().Values[1].(string)})
		}
		return contacts, nil

	})
	if err != nil {
		return nil, err
	}
	contactResponse := &pb.ContactsResponse{Contacts: contacts.([]*pb.Contact)}
	return contactResponse, nil

}
