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

func NewChat(id, userIDa, userIDb string) Chat {
	return Chat{Id: converter.GetObjectId(id), UserIDa: userIDa, UserIDb: userIDb}

}
