package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var chats = []*domain.Chat{
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e2"), // rasti djordje
		UserIDa:       "62752bf27407f54ce1839cb9",
		UserIDb:       "62752bf27407f54ce1839cb8",
		UserASeenDate: time.Date(2022, time.June, 11, 16, 52, 14, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 12, 11, 0, 0, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Cao kako si sta radis?",
				Date:         time.Date(2022, time.June, 11, 16, 34, 12, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "Dobro sam, evo radim xws :)",
				Date:         time.Date(2022, time.June, 11, 16, 51, 14, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "Sta ti radis?",
				Date:         time.Date(2022, time.June, 11, 16, 52, 53, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "Danas je jako lepo vreme",
				Date:         time.Date(2022, time.June, 11, 17, 5, 14, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e3"), // rasti nikola luburic
		UserIDa:       "62752bf27407f54ce1839cb9",
		UserIDb:       "62752bf27407f54ce1839cb4",
		UserASeenDate: time.Date(2022, time.June, 12, 10, 0, 0, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 12, 11, 0, 0, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Pozdrav profesore kada ce biti ispit iz PSW? :)",
				Date:         time.Date(2022, time.June, 5, 12, 34, 12, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb4",
				Text:         "Bice u ispitnom roku",
				Date:         time.Date(2022, time.June, 5, 13, 35, 14, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb4",
				Text:         "Pratite studentsku sluzbu, nadam se da vam je moj odgovor bio koristan",
				Date:         time.Date(2022, time.June, 5, 13, 36, 46, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Jako",
				Date:         time.Date(2022, time.June, 6, 10, 2, 12, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e5"),
		UserIDa:       "62752bf27407f54ce1839cb7", //tara
		UserIDb:       "62752bf27407f54ce1839cb9", //rasti
		UserASeenDate: time.Date(2022, time.June, 13, 19, 11, 20, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 13, 19, 10, 5, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Kako ide diplomski? :)",
				Date:         time.Date(2022, time.June, 13, 19, 10, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb7",
				Text:         "Odlicno uradila sam ga jos pre 3 meseca, sada ispravljam rad menjam iz latinice u cirilicu",
				Date:         time.Date(2022, time.June, 13, 19, 11, 15, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb7",
				Text:         "Neke reci se ne mogu tako lako prevesti na cirilicu, haha",
				Date:         time.Date(2022, time.June, 13, 19, 11, 20, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb7",
				Text:         "Jesi li ti poceo diplomski da radis?",
				Date:         time.Date(2022, time.June, 13, 19, 11, 20, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e4"),
		UserIDa:       "62752bf27407f54ce1839cb9", //rasti
		UserIDb:       "62752bf27407f54ce1839cb6", // zarko
		UserASeenDate: time.Date(2022, time.June, 14, 12, 47, 0, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 14, 12, 47, 0, 0, time.UTC),
		Messages:      []domain.Message{},
	},
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func getIdFromHex(userID string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(userID)
	return id
}
