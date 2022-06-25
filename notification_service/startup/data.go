package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var notifications = []*domain.Notification{
	//TODO: fill with data
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
