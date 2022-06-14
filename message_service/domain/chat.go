package domain

import (
	converter "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/converter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Chat converter.GetObjectId(id)
type Chat struct {
	Id            primitive.ObjectID `bson:"_id"`
	UserIDa       string             `bson:"userIDa"`
	UserIDb       string             `bson:"userIDb"`
	UserASeenDate time.Time          `bson:"userASeenDate"`
	UserBSeenDate time.Time          `bson:"userBSeenDate"`
	Messages      []Message          `bson:"messages"`
}

func (c Chat) GetSeenDateByUserID(userID string) time.Time {
	if userID == c.UserIDa {
		return c.UserASeenDate
	} else {
		return c.UserBSeenDate
	}
}

func (c Chat) HaveUserID(userID string) bool {
	return c.UserIDa == userID || c.UserIDb == userID
}

func (c Chat) GetOtherUserID(myUserID string) string {
	if c.UserIDa == myUserID {
		return c.UserIDb
	} else if c.UserIDb == myUserID {
		return c.UserIDa
	}
	return c.UserIDa
}

func (c Chat) UpdateAuthorSeenDate(authorID string, time time.Time) {
	if c.UserIDa == authorID {
		c.UserASeenDate = time
	} else if c.UserIDb == authorID {
		c.UserBSeenDate = time
	}
}

func NewChat(id, userIDa, userIDb string) Chat {
	return Chat{Id: converter.GetObjectId(id), UserIDa: userIDa, UserIDb: userIDb}

}
