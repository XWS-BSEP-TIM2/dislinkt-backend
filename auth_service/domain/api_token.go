package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ApiToken struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id"`
	ApiCode      string             `bson:"token_code"`
	CreationDate time.Time          `bson:"creation_date"`
}
