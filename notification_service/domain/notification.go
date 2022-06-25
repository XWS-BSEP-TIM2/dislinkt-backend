package domain

import (
	"time"
)

type Notification struct {
	OwnerId    string    `bson:"ownerId"`
	ForwardUrl string    `bson:"forwardUrl"`
	Text       string    `bson:"text"`
	Date       time.Time `bson:"date"`
}
