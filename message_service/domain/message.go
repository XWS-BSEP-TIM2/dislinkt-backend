package domain

import (
	"time"
)

type Message struct {
	AuthorUserID string    `bson:"authorUserID"`
	Text         string    `bson:"text"`
	Date         time.Time `bson:"date"`
}
