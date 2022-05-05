package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var users = []*domain.User{
	{
		Id:       getObjectId("ID000001"),
		Username: "marko99",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000002"),
		Username: "srdjan",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000003"),
		Username: "nikola93",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000004"),
		Username: "svetozar",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000005"),
		Username: "zarkoo",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000006"),
		Username: "tara",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000007"),
		Username: "djordje",
		Password: "pass",
		Role:     domain.USER,
	},
	{
		Id:       getObjectId("ID000008"),
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
