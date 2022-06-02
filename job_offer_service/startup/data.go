package startup

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jobOffers = []*domain.JobOffer{
	//{
	//	Id:           getIdFromHex("62752bf27407f54ce1839cb2"),
	//	Position:     "Marko",
	//	Seniority:    "Markovic",
	//	Description:  "marko99",
	//	Technologies: nil,
	//},
	//{
	//	Id:          getIdFromHex("62752bf27407f54ce1839cb3"),
	//	Position:    "Srdjan",
	//	Seniority:   "Srdjanovic",
	//	Description: "srdjan",
	//
	//	Technologies: nil,
	//},
	//{
	//	Id:          getIdFromHex("62752bf27407f54ce1839cb4"),
	//	Position:    "Nikola",
	//	Seniority:   "Luburic",
	//	Description: "nikola93",
	//
	//	Technologies: nil,
	//},
	//{
	//	Id:          getIdFromHex("62752bf27407f54ce1839cb5"),
	//	Position:    "Sveto",
	//	Seniority:   "Svetozar",
	//	Description: "svetozar",
	//
	//	Technologies: nil,
	//},
	//{
	//	Id:           getIdFromHex("62752bf27407f54ce1839cb6"),
	//	Position:     "Zarko",
	//	Seniority:    "Blagojevic",
	//	Description:  "zarkoo",
	//	Technologies: nil,
	//},
	//{
	//	Id:           getIdFromHex("62752bf27407f54ce1839cb7"),
	//	Position:     "Tara",
	//	Seniority:    "Pogancev",
	//	Description:  "Jelovceva",
	//	Technologies: nil,
	//},
	//{
	//	Id:           getIdFromHex("62752bf27407f54ce1839cb8"),
	//	Position:     "Djordje",
	//	Seniority:    "Krsmanovic",
	//	Description:  "djordje",
	//	Technologies: nil,
	//},
	//{
	//	Id:           getIdFromHex("62752bf27407f54ce1839cb9"),
	//	Position:     "Rastislav",
	//	Seniority:    "Kukucka",
	//	Description:  "rasti",
	//	Technologies: nil,
	//},
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
