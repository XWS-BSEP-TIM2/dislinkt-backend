package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSettings struct {
	Id                      primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId                 primitive.ObjectID `bson:"ownerId"`
	PostNotifications       bool               `bson:"postNotifications"`
	ConnectionNotifications bool               `bson:"connectionNotifications"`
	MessageNotifications    bool               `bson:"messageNotifications"`
}
