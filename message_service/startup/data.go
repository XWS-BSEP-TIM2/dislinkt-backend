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
				Text:         "Dobro sam, evo radim xws :) inace uspeo sam da uradim onaj deo u vezi 2 autentifikacije, bas je bilo zanimljivo i naucio sam bas cool stvari",
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
		UserBSeenDate: time.Date(2022, time.June, 5, 13, 40, 47, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Pozdrav profesore kada ce biti ispit iz PSW? :) Takodje imam pitanje u vezi TDD da li mora uvek da se pise Test pre implementacije ili ne mora. npr mi sada radimo xws i tu u Go nismo jos radili testove, da li preporucujete da uradimo neke testove ili ne mora",
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
				AuthorUserID: "62752bf27407f54ce1839cb4",
				Text:         "Sto se tice Test Driven Developmenta, preporucujem da pocnete da pisete testove u Go, nije toliko tesko postoji biblioteka bas za to namenjena zove se testing, ako treba pomoc pisite ja sam vec testirao sa njom",
				Date:         time.Date(2022, time.June, 5, 13, 37, 34, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Jako",
				Date:         time.Date(2022, time.June, 6, 10, 2, 12, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Hvala vam puno na odgovoru",
				Date:         time.Date(2022, time.June, 6, 10, 2, 15, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Vidimo se na ispitu, nadam se 10 ;)",
				Date:         time.Date(2022, time.June, 6, 10, 2, 25, 0, time.UTC),
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
		UserASeenDate: time.Date(2022, time.June, 26, 12, 55, 20, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 26, 12, 49, 0, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Kako je na praksi iz AI sta radite sada, jel ste vec zavrsili chat bota? :)",
				Date:         time.Date(2022, time.June, 26, 12, 47, 0, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb6",
				Text:         "Hvala na pitanju odlicno je, stvarno je cool, uspeo sam da napravim jako korisnog chet bot, i pokrio jako puno slucaja :D",
				Date:         time.Date(2022, time.June, 26, 12, 49, 0, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Svaka cast, Odlicno jako mi je drago sto ti se svidja",
				Date:         time.Date(2022, time.June, 26, 12, 51, 20, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "I meni je to bas cool, nlp :) mada openCV je strava",
				Date:         time.Date(2022, time.June, 26, 12, 55, 20, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e6"),
		UserIDa:       "62752bf27407f54ce1839cb9", //rasti
		UserIDb:       "62752bf27407f54ce1839cb3", //srdjan
		UserASeenDate: time.Date(2022, time.June, 26, 12, 12, 4, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 26, 12, 14, 4, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb3",
				Text:         "Sta radis ovih dana jel hoces da idemo na bazen?",
				Date:         time.Date(2022, time.June, 26, 12, 10, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Aj tipa posle 12 jula kad nam zavrse obaveze",
				Date:         time.Date(2022, time.June, 26, 12, 12, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb3",
				Text:         "Moze bice cool zovemo ekipu sa faxa",
				Date:         time.Date(2022, time.June, 26, 12, 14, 4, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e7"),
		UserIDa:       "62752bf27407f54ce1839cb9", //rasti
		UserIDb:       "62752bf27407f54ce1839cb2", //marko
		UserASeenDate: time.Date(2022, time.June, 25, 8, 10, 4, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 25, 9, 10, 4, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb9",
				Text:         "Da li si kupio knjige za srednju skolu? Moj brat prodaje :)",
				Date:         time.Date(2022, time.June, 25, 8, 10, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb2",
				Text:         "Ee kupio sam hvala, nabavio sam vec",
				Date:         time.Date(2022, time.June, 25, 9, 10, 4, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e8"),
		UserIDa:       "62752bf27407f54ce1839cb9", // rasti
		UserIDb:       "62752bf27407f54ce1839cb5", //svetozar
		UserASeenDate: time.Date(2022, time.June, 10, 20, 11, 15, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 11, 21, 10, 7, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb5",
				Text:         "Dislinkt je jako lepa drustvena mreza",
				Date:         time.Date(2022, time.June, 24, 17, 10, 4, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23e9"),
		UserIDa:       "62752bf27407f54ce1839cb6", // zarko
		UserIDb:       "62752bf27407f54ce1839cb7", //tara
		UserASeenDate: time.Date(2022, time.June, 21, 21, 15, 4, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 22, 21, 13, 4, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb6",
				Text:         "Kada nam ono bese dolazi ispit iz inteligentnih?",
				Date:         time.Date(2022, time.June, 21, 21, 13, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb6",
				Text:         "Mislim da stizemo lagano",
				Date:         time.Date(2022, time.June, 21, 21, 15, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb7",
				Text:         "3 Jula :) ",
				Date:         time.Date(2022, time.June, 22, 21, 13, 4, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23ea"),
		UserIDa:       "62752bf27407f54ce1839cb6", // zarko
		UserIDb:       "62752bf27407f54ce1839cb8", //djordje
		UserASeenDate: time.Date(2022, time.June, 22, 21, 16, 14, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 22, 21, 13, 14, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "Kako je bilo na keju preksinoc?",
				Date:         time.Date(2022, time.June, 22, 21, 13, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "Jel bilo drustva?",
				Date:         time.Date(2022, time.June, 22, 21, 13, 14, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb6",
				Text:         "uu bilo je odlicno, bas je bilo dobro druzenje izblejali smo, aj dodji veceras opet smo tu",
				Date:         time.Date(2022, time.June, 22, 21, 16, 14, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23eb"),
		UserIDa:       "62752bf27407f54ce1839cb6", // zarko
		UserIDb:       "62752bf27407f54ce1839cb5", //svetozar
		UserASeenDate: time.Date(2022, time.June, 10, 23, 11, 20, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 11, 21, 12, 53, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb5",
				Text:         "Pozdrav :)",
				Date:         time.Date(2022, time.June, 11, 21, 12, 53, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23ec"),
		UserIDa:       "62752bf27407f54ce1839cb7", // tara
		UserIDb:       "62752bf27407f54ce1839cb8", //djordje
		UserASeenDate: time.Date(2022, time.June, 23, 14, 31, 14, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 23, 14, 31, 14, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb7",
				Text:         "Kako ti se cini novi Apple Macbook Air 2022?",
				Date:         time.Date(2022, time.June, 23, 14, 22, 4, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "A ne znam, posalji mi specifikacije da vidim",
				Date:         time.Date(2022, time.June, 23, 14, 23, 14, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb7",
				Text:         "Apple M2 chip\n\n8-core CPU with 4 performance cores and 4 efficiency cores\n8-core GPU\n16-core Neural Engine\n100GB/s memory bandwidth\nMedia engine\n\nHardware-accelerated H.264, HEVC, ProRes, and ProRes RAW\nVideo decode engine\nVideo encode engine\nProRes encode and decode engine\nConfigurable to:\n\nM2 with 8-core CPU and 10-core GPU",
				Date:         time.Date(2022, time.June, 23, 14, 25, 14, 0, time.UTC),
			},
			{
				AuthorUserID: "62752bf27407f54ce1839cb8",
				Text:         "Uuu Mocan bas",
				Date:         time.Date(2022, time.June, 23, 14, 31, 14, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23ed"),
		UserIDa:       "62752bf27407f54ce1839cb7", // tara
		UserIDb:       "62752bf27407f54ce1839cb5", //svetozar
		UserASeenDate: time.Date(2022, time.June, 10, 10, 11, 24, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 11, 10, 13, 5, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb5",
				Text:         "Hello :)",
				Date:         time.Date(2022, time.June, 11, 10, 13, 5, 0, time.UTC),
			},
		},
	},
	{
		Id:            getIdFromHex("62a76abf5b14e448f4bd23ee"),
		UserIDa:       "62752bf27407f54ce1839cb8", // djordje
		UserIDb:       "62752bf27407f54ce1839cb3", //srdjan
		UserASeenDate: time.Date(2022, time.June, 10, 25, 11, 20, 0, time.UTC),
		UserBSeenDate: time.Date(2022, time.June, 11, 12, 10, 5, 0, time.UTC),
		Messages: []domain.Message{
			{
				AuthorUserID: "62752bf27407f54ce1839cb3",
				Text:         "De si Djordje brate :D",
				Date:         time.Date(2022, time.June, 11, 12, 10, 5, 0, time.UTC),
			},
		},
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
