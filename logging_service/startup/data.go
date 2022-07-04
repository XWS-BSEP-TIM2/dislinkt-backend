package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/domain"
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

func getIdFromHex(objectId string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(objectId)
	return id
}

var events = []*domain.Event{
	{
		UserId:      userIdMap[tara],
		Title:       "Registration",
		Description: "A new user has registered",
		Date:        time.Date(2022, time.June, 20, 10, 20, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[sveta],
		Title:       "Registration",
		Description: "A new user has registered",
		Date:        time.Date(2022, time.May, 22, 11, 50, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[djole],
		Title:       "Registration",
		Description: "A new user has registered",
		Date:        time.Date(2022, time.February, 23, 12, 45, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[zare],
		Title:       "Registration",
		Description: "A new user has registered",
		Date:        time.Date(2022, time.January, 9, 13, 30, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[rasti],
		Title:       "Registration",
		Description: "A new user has registered",
		Date:        time.Date(2022, time.February, 10, 14, 16, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[rasti],
		Title:       "Post",
		Description: "User has published a new post",
		Date:        time.Date(2022, time.April, 22, 20, 52, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[zare],
		Title:       "Post",
		Description: "User has published a new post",
		Date:        time.Date(2022, time.February, 7, 21, 49, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[djole],
		Title:       "Post",
		Description: "User has published a new post",
		Date:        time.Date(2022, time.January, 6, 22, 45, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[tara],
		Title:       "Job Offer",
		Description: "User has made a new job offer",
		Date:        time.Date(2022, time.March, 4, 23, 0, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[tara],
		Title:       "Job Offer",
		Description: "User has made a new job offer",
		Date:        time.Date(2022, time.March, 2, 8, 12, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[zare],
		Title:       "Job Offer",
		Description: "User has made a new job offer",
		Date:        time.Date(2022, time.May, 26, 7, 36, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[zare],
		Title:       "Job Offer",
		Description: "User has made a new job offer",
		Date:        time.Date(2022, time.June, 24, 6, 45, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[zare],
		Title:       "Update Profile",
		Description: "User has updated their profile settings",
		Date:        time.Date(2022, time.March, 16, 2, 40, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[zare],
		Title:       "Job Offer Update",
		Description: "User has updated their job offer",
		Date:        time.Date(2022, time.February, 14, 1, 30, 0, 10000000, time.UTC),
	},
	{
		UserId:      userIdMap[tara],
		Title:       "Job Offer Delete",
		Description: "User has deleted their job offer",
		Date:        time.Date(2022, time.January, 12, 12, 20, 0, 10000000, time.UTC),
	},
}
