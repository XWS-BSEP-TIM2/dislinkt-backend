package domain

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Experience struct {
	Id               primitive.ObjectID   `bson:"_id"`
	Name             string               `bson:"name"`
	TypeOfExperience enums.ExperienceType `bson:"type_of_experience"`
	Description      string               `bson:"description"`
	StartDate        time.Time            `bson:"start_date"`
	EndDate          time.Time            `bson:"end_date"`
}

func NewExperience(name string, expType enums.ExperienceType, description string, startdate time.Time, endDate time.Time) Experience {
	return Experience{Name: name, TypeOfExperience: expType, Description: description, StartDate: startdate, EndDate: endDate}
}
