package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var users = []*domain.Profile{
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb2"),
		Name:        "Marko",
		Surname:     "Markovic",
		Username:    "marko99",
		Email:       "marko99@gmail.com",
		Biography:   "Ma ja sam marko",
		Gender:      0,
		PhoneNumber: "062978354",
		BirthDate:   time.Date(2009, time.November, 10, 10, 0, 0, 0, time.UTC),
		IsPrivate:   false,
		Skills:      []domain.Skill{},
		Experiences: []domain.Experience{},
		IsTwoFactor: false,
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb3"),
		Name:        "Srdjan",
		Surname:     "Srdjanovic",
		Username:    "srdjan",
		Email:       "srdjan@gmail.com",
		Biography:   "Srdajn neki lik super",
		Gender:      0,
		PhoneNumber: "06234564",
		BirthDate:   time.Date(1990, time.January, 14, 10, 0, 0, 0, time.UTC),
		IsPrivate:   true,
		Skills:      []domain.Skill{},
		Experiences: []domain.Experience{},
		IsTwoFactor: false,
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb4"),
		Name:        "Nikola",
		Surname:     "Luburic",
		Username:    "nikola93",
		Email:       "nikola@gmail.com",
		Biography:   "Pro Doktor",
		Gender:      0,
		PhoneNumber: "063589625",
		BirthDate:   time.Date(1993, time.January, 14, 10, 0, 0, 0, time.UTC),
		IsPrivate:   false,
		Skills:      []domain.Skill{},
		Experiences: []domain.Experience{},
		IsTwoFactor: false,
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb5"),
		Name:        "Sveto",
		Surname:     "Svetozar",
		Username:    "svetozar",
		Email:       "svet@gmail.com",
		Biography:   "Sve sve svetozar :)",
		Gender:      0,
		PhoneNumber: "0634684853",
		BirthDate:   time.Date(1998, time.January, 14, 10, 0, 0, 0, time.UTC),
		IsPrivate:   true,
		Skills:      []domain.Skill{},
		Experiences: []domain.Experience{},
		IsTwoFactor: false,
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb6"),
		Name:        "Zarko",
		Surname:     "Blagojevic",
		Username:    "zarkoo",
		Email:       "zare00@gmail.com",
		Biography:   "Student at Faculty of Technical Sciences, University of Novi Sad",
		Gender:      0,
		PhoneNumber: "063687626",
		BirthDate:   time.Date(2000, time.February, 18, 10, 0, 0, 0, time.UTC),
		IsPrivate:   false,
		Skills: []domain.Skill{
			{
				Id:   getObjectId("IDS00001"),
				Name: "Tensorflow",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00002"),
				Name: "Docker",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00003"),
				Name: "CircleCI",
				Type: 0,
			},
		},
		Experiences: []domain.Experience{
			{
				Id:               getObjectId("IDE000001"),
				Name:             "FTN - Computer Science",
				Description:      "Applied Computer Science and Informatics student, at the Faculty of Technical Sciences",
				TypeOfExperience: 0,
				StartDate:        time.Date(2018, time.October, 1, 10, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2022, time.June, 30, 10, 0, 0, 0, time.UTC),
			},
		},
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb7"),
		Name:        "Tara",
		Surname:     "Pogancev",
		Username:    "Jelovceva",
		Email:       "tara00@gmail.com",
		Biography:   "Applied Computer Science and Informatics student, at the Faculty of Technical Sciences, University of Novi Sad. Highly ambitious, open towards new professional experiences and skills. Distinctly interested in acquiring experience in software and web development",
		Gender:      1,
		PhoneNumber: "063687626",
		BirthDate:   time.Date(2000, time.February, 3, 10, 0, 0, 0, time.UTC),
		IsPrivate:   false,
		Skills: []domain.Skill{
			{
				Id:   getObjectId("IDS00004"),
				Name: "Java",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00005"),
				Name: "Angular",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00006"),
				Name: "Bulma",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00007"),
				Name: "TypeScript",
				Type: 0,
			},
		},
		Experiences: []domain.Experience{
			{
				Id:               getObjectId("IDE000001"),
				Name:             "FTN - Computer Science",
				Description:      "Applied Computer Science and Informatics student, at the Faculty of Technical Sciences",
				TypeOfExperience: 0,
				StartDate:        time.Date(2018, time.October, 1, 10, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2022, time.June, 30, 10, 0, 0, 0, time.UTC),
			},
		},
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb8"),
		Name:        "Djordje",
		Surname:     "Krsmanovic",
		Username:    "djordje",
		Email:       "djordje99@gmail.com",
		Biography:   "Student at Faculty of Technical Sciences, University of Novi Sad",
		Gender:      0,
		PhoneNumber: "06284684",
		BirthDate:   time.Date(1999, time.January, 14, 10, 0, 0, 0, time.UTC),
		IsPrivate:   false,
		Skills: []domain.Skill{
			{
				Id:   getObjectId("IDS00008"),
				Name: "Java",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00009"),
				Name: "GO",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00010"),
				Name: "Docker",
				Type: 0,
			},
		},
		Experiences: []domain.Experience{
			{
				Id:               getObjectId("IDE000001"),
				Name:             "FTN - Computer Science",
				Description:      "Applied Computer Science and Informatics student, at the Faculty of Technical Sciences",
				TypeOfExperience: 0,
				StartDate:        time.Date(2018, time.October, 1, 10, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2022, time.June, 30, 10, 0, 0, 0, time.UTC),
			},
		},
		IsTwoFactor: true,
	},
	{
		Id:          getIdFromHex("62752bf27407f54ce1839cb9"),
		Name:        "Rastislav",
		Surname:     "Kukucka",
		Username:    "rasti",
		Email:       "rasti@gmail.com",
		Biography:   "Student at Faculty of Technical Sciences, University of Novi Sad",
		Gender:      0,
		PhoneNumber: "0627458242",
		BirthDate:   time.Date(2000, time.January, 7, 10, 0, 0, 0, time.UTC),
		IsPrivate:   false,
		Skills: []domain.Skill{
			{
				Id:   getObjectId("IDS00011"),
				Name: "Python",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00012"),
				Name: "Java",
				Type: 0,
			},
			{
				Id:   getObjectId("IDS00013"),
				Name: "Angular",
				Type: 0,
			},
		},
		Experiences: []domain.Experience{
			{
				Id:               getObjectId("IDE000001"),
				Name:             "FTN - Computer Science",
				Description:      "Applied Computer Science and Informatics student, at the Faculty of Technical Sciences",
				TypeOfExperience: 0,
				StartDate:        time.Date(2018, time.October, 1, 10, 0, 0, 0, time.UTC),
				EndDate:          time.Date(2022, time.June, 30, 10, 0, 0, 0, time.UTC),
			},
		},
		IsTwoFactor: true,
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
