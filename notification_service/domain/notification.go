package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	OwnerId      primitive.ObjectID `bson:"ownerId"`
	ForwardUrl   string             `bson:"forwardUrl"`
	Text         string             `bson:"text"`
	Date         time.Time          `bson:"date"`
	Seen         bool               `bson:"seen"`
	UserFullName string             `bson:"userFullName"`
}
