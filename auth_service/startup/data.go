package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var users = []*domain.User{
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb2"),
		Username: "marko99",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb3"),
		Username: "srdjan",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb4"),
		Username: "nikola93",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb5"),
		Username: "svetozar",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb6"),
		Username: "zarkoo",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb7"),
		Username: "Jelovceva",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb8"),
		Username: "djordje",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb9"),
		Username: "rasti",
		Password: "pass",
		Role:     domain.USER,
	},
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
