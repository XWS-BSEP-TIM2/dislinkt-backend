package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PasswordlessToken struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       primitive.ObjectID `bson:"user_id"`
	TokenCode    string             `bson:"token_code"`
	CreationDate time.Time          `bson:"creation_date"`
}
