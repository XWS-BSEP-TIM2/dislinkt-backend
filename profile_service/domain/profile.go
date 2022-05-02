package domain

import (
	converter "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/converter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Profile struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Surname     string             `bson:"surname"`
	Username    string             `bson:"username"`
	Email       string             `bson:"email"`
	Biography   string             `bson:"biography"`
	Gender      enums.Gender       `bson:"gender"`
	PhoneNumber string             `bson:"phoneNumber"`
	BirthDate   time.Time          `bson:"birthDate"`
	IsPrivate   bool               `bson:"isPrivate"`
	Skills      []Skill            `bson:"skills"`
	Experiences []Experience       `bson:"experiences"`
}

func NewProfile(id string, name string, surname string, username string, email string, biography string, gender enums.Gender, phoneNumber string, birthDate time.Time, isPrivate bool) Profile {
	return Profile{Id: converter.GetObjectId(id), Name: name, Surname: surname, Username: username, Email: email, Biography: biography, Gender: gender, PhoneNumber: phoneNumber, BirthDate: birthDate, IsPrivate: isPrivate}

}
