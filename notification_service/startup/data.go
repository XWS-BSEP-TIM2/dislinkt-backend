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

var userSettings = []*domain.UserSettings{
	{
		OwnerId:                 userIdMap[tara],
		PostNotifications:       true,
		ConnectionNotifications: true,
		MessageNotifications:    false,
	},
}

var notifications = []*domain.Notification{
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ac4"),
		OwnerId:      userIdMap[tara],
		ForwardUrl:   "posts/62b75d956a06acad63b03c33",
		Text:         "posted on their profile",
		Date:         time.Date(2022, time.June, 22, 10, 0, 0, 10000000, time.UTC),
		Seen:         false,
		UserFullName: "Sveto Svetozar",
	},
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ac5"),
		OwnerId:      userIdMap[tara],
		ForwardUrl:   "https://localhost:4200/posts/62b75d956a06acad63b03c27",
		Text:         "posted on their profile",
		Date:         time.Date(2022, time.June, 20, 10, 0, 0, 10000000, time.UTC),
		Seen:         false,
		UserFullName: "Zarko Blagojevic",
	},
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ac6"),
		OwnerId:      userIdMap[tara],
		ForwardUrl:   "chat",
		Text:         "sent you a message",
		Date:         time.Date(2022, time.March, 19, 16, 30, 0, 10000000, time.UTC),
		Seen:         true,
		UserFullName: "Rastislav Kukucka",
	},
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ac7"),
		OwnerId:      userIdMap[rasti],
		ForwardUrl:   "posts/62b75d956a06acad63b03c37",
		Text:         "posted on their profile",
		Date:         time.Date(2022, time.February, 18, 10, 48, 0, 10000000, time.UTC),
		Seen:         true,
		UserFullName: "Djordje Krsmanovic",
	},
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ac8"),
		OwnerId:      userIdMap[zare],
		ForwardUrl:   "posts/62b75d956a06acad63b03c37",
		Text:         "posted on their profile",
		Date:         time.Date(2022, time.February, 18, 10, 48, 0, 10000000, time.UTC),
		Seen:         true,
		UserFullName: "Djordje Krsmanovic",
	},
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ac9"),
		OwnerId:      userIdMap[tara],
		ForwardUrl:   "posts/62b75d956a06acad63b03c37",
		Text:         "posted on their profile",
		Date:         time.Date(2022, time.February, 18, 10, 48, 0, 10000000, time.UTC),
		Seen:         true,
		UserFullName: "Djordje Krsmanovic",
	},
	{
		Id:           getIdFromHex("62b739b2c74d379ebaf20ad0"),
		OwnerId:      userIdMap[tara],
		ForwardUrl:   "profile/62752bf27407f54ce1839cb7/requests",
		Text:         "sent you a friend request",
		Date:         time.Date(2022, time.January, 3, 10, 0, 0, 10000000, time.UTC),
		Seen:         true,
		UserFullName: "Djordje Krsmanovic",
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
