package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var users = []*domain.User{
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb2"),
		Username: "marko99",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "marko99@gmail.com",
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb3"),
		Username: "srdjan",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "srdjan@gmail.com",
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb4"),
		Username: "nikola93",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "nikola@gmail.com",
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb5"),
		Username: "svetozar",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "svet@gmail.com",
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb6"),
		Username: "zarkoo",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "zare00@gmail.com",
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb7"),
		Username: "Jelovceva",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "tara00@gmail.com",
	},
	{
		Id:           getIdFromHex("62752bf27407f54ce1839cb8"),
		Username:     "djordje",
		Password:     "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:         domain.USER,
		Verified:     true,
		Email:        "djordje1499@gmail.com",
		IsTFAEnabled: true,
	},
	{
		Id:       getIdFromHex("62752bf27407f54ce1839cb9"),
		Username: "rasti",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.USER,
		Verified: true,
		Email:    "dislinktx@gmail.com",
	},
	{
		Id:       getIdFromHex("62752bf27407f85de1839cb9"),
		Username: "admin",
		Password: "$2a$05$ic587IZBvYNLqSRSsTuae.nZISenVmwmg0ddr8DmSI11lZV6VWixu",
		Role:     domain.ADMIN,
		Verified: true,
		Email:    "dislinktx@gmail.com",
	},
}

var apiTokens = []*domain.ApiToken{
	{
		Id:           getIdFromHex("62752bf27407f51ce1839cb0"),
		UserId:       getIdFromHex("62752bf27407f54ce1839cb8"),
		CreationDate: time.Now(),
		ApiCode:      "adadqbek123krmtgk123e1rfd",
	},
	{
		Id:           getIdFromHex("62752bf27407f51ce1839cb1"),
		UserId:       getIdFromHex("62752bf27407f54ce1839cb7"),
		CreationDate: time.Now(),
		ApiCode:      "adadqbek123krmtgk123e1rff",
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
