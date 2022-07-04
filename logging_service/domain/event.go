package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	UserId      primitive.ObjectID `bson:"userId"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Date        time.Time          `bson:"date"`
}
