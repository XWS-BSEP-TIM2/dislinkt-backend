package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const ( // iota is reset to 0
	tara  = iota
	sveta = iota
	djole = iota
	rasti = iota
	zare  = iota
)

var userIdMap = map[int]primitive.ObjectID{
	tara:  getIdFromHex("62752bf27407f54ce1839cb7"),
	sveta: getIdFromHex("62752bf27407f54ce1839cb5"),
	djole: getIdFromHex("62752bf27407f54ce1839cb8"),
	rasti: getIdFromHex("62752bf27407f54ce1839cb9"),
	zare:  getIdFromHex("62752bf27407f54ce1839cb6"),
}

var notifications = []*domain.Notification{
	{
		Id:         getIdFromHex("62b739b2c74d379ebaf20ac4"),
		OwnerId:    userIdMap[tara],
		ForwardUrl: "ulr",
		Text:       "User Zarko Blagojevic posted on their profile",
		Date:       time.Date(2022, time.January, 5, 10, 0, 0, 10000000, time.UTC),
		Seen:       true,
	},
	{
		Id:         getIdFromHex("62b739b2c74d379ebaf20ac5"),
		OwnerId:    userIdMap[tara],
		ForwardUrl: "ulr",
		Text:       "User Rastislav Kukucka sent you a message",
		Date:       time.Date(2022, time.January, 5, 16, 30, 0, 10000000, time.UTC),
		Seen:       true,
	},
	{
		Id:         getIdFromHex("62b739b2c74d379ebaf20ac6"),
		OwnerId:    userIdMap[rasti],
		ForwardUrl: "ulr",
		Text:       "User Zarko Blagojevic posted on their profile",
		Date:       time.Date(2022, time.April, 18, 10, 48, 0, 10000000, time.UTC),
		Seen:       true,
	},
	{
		Id:         getIdFromHex("62b739b2c74d379ebaf20ac7"),
		OwnerId:    userIdMap[tara],
		ForwardUrl: "ulr",
		Text:       "User Djordje Krsmanovic sent you a connection request",
		Date:       time.Date(2022, time.June, 5, 10, 0, 0, 10000000, time.UTC),
		Seen:       false,
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func getIdFromHex(objectId string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(objectId)
	return id
}
