package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobOffer struct {
	Id                 primitive.ObjectID `bson:"_id"`
	UserId             primitive.ObjectID `bson:"user_id"`
	JobOfferUniqueCode string             `bson:"job_offer_unique_code"`
	CompanyName        string             `bson:"company_name"`
	Position           string             `bson:"position"`
	Seniority          string             `bson:"seniority"`
	Description        string             `bson:"description"`
	Technologies       []string           `bson:"technologies"`
}
