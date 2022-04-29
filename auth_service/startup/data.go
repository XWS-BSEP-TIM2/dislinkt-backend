package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var products = []*domain.User{
	{
		Id:       getObjectId("623b0cc3a34d25d8567f9f82"),
		Username: "name",
		Password: "123",
		Role:     domain.USER,
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}
