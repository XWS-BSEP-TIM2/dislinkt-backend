package persistence

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"time"
)

func checkIfUserExist(userID string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (existing_uer:USER) WHERE existing_uer.userID = $userID RETURN existing_uer.userID",
		map[string]interface{}{"userID": userID})

	if result != nil && result.Next() && result.Record().Values[0] == userID {
		return true
	}
	return false
}

func checkIfFriendExist(userIDa, userIDb string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
			"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
			"MATCH (u1)-[r:FRIEND]->(u2) "+
			"RETURN r.date ",
		map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb})

	if result != nil && result.Next() {
		return true
	}
	return false
}

func checkIfBlockExist(userIDa, userIDb string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
			"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
			"MATCH (u1)-[r:BLOCK]->(u2) "+
			"RETURN r.date ",
		map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb})

	if result != nil && result.Next() {
		return true
	}
	return false
}

func blockUser(userIDa, userIDb string, transaction neo4j.Transaction) bool {

	dateNow := time.Now().Local().Unix()
	result, err := transaction.Run(
		"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
			"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
			"CREATE (u1)-[r:BLOCK {date: $dateNow}]->(u2) "+
			"RETURN r.date",
		map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb, "dateNow": dateNow, "msgID": "newMsgID"})

	if err != nil {
		fmt.Println(err)
		return false
	}
	if result != nil && result.Next() {
		return true
	}
	return false
}

func removeFriend(userIDa, userIDb string, transaction neo4j.Transaction) bool {
	result, err := transaction.Run(
		"MATCH (u1:USER{userID: $uIDa})-[r:FRIEND]->(u2:USER{userID: $uIDb}) DELETE r RETURN u1.userID",
		map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb})

	if err != nil {
		fmt.Println(err)
		return false
	}
	if result != nil && result.Next() {
		return true
	}
	return false
}

func unblockUser(userIDa, userIDb string, transaction neo4j.Transaction) bool {
	result, err := transaction.Run(
		"MATCH (u1:USER{userID: $uIDa})-[r:BLOCK]->(u2:USER{userID: $uIDb}) DELETE r RETURN u1.userID",
		map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb})

	if err != nil {
		fmt.Println(err)
		return false
	}
	if result != nil && result.Next() {
		return true
	}
	return false
}

func getFriendsOfFriendsButNotBlockedRecommendation(userID string, transaction neo4j.Transaction) ([]*domain.UserConn, error) {
	result, err := transaction.Run(
		"MATCH (u1:USER)-[:FRIEND]->(u2:USER)<-[:FRIEND]-(u3:USER) "+
			"WHERE u1.userID=$uID AND u3.userID<>$uID "+
			"AND NOT exists((u1)-[:FRIEND]-(u3)) "+
			"AND NOT exists((u1)-[:BLOCK]-(u3)) "+
			"RETURN distinct u3.userID, u3.isPublic "+
			"LIMIT 20 ",
		map[string]interface{}{"uID": userID})

	if err != nil {
		return nil, err
	}

	var recommendation []*domain.UserConn
	for result.Next() {
		recommendation = append(recommendation, &domain.UserConn{UserID: result.Record().Values[0].(string), IsPublic: result.Record().Values[1].(bool)})
	}
	return recommendation, nil
}

func getFriendRecommendation(userID string, transaction neo4j.Transaction) ([]*domain.UserConn, error) {
	result, err := transaction.Run(
		"MATCH (u1:USER) "+
			"MATCH (u2:USER)-[r:FRIEND]->(:USER) "+
			"WHERE u1.userID=$uID AND u2.userID<>$uID "+
			"AND NOT exists((u1)-[:FRIEND]-(u2)) "+
			"AND NOT exists((u1)-[:BLOCK]-(u2)) "+
			"RETURN distinct u2.userID, u2.isPublic, COUNT(r) as num_of_friends "+
			"ORDER BY num_of_friends DESC "+
			"LIMIT 20 ",
		map[string]interface{}{"uID": userID})

	if err != nil {
		return nil, err
	}

	var recommendation []*domain.UserConn
	for result.Next() {
		recommendation = append(recommendation, &domain.UserConn{UserID: result.Record().Values[0].(string), IsPublic: result.Record().Values[1].(bool)})
	}
	return recommendation, nil
}
