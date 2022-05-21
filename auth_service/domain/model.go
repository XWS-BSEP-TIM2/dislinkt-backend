package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Role int

const (
	USER Role = iota
	ADMIN
)

type User struct {
	Id                       primitive.ObjectID `bson:"_id,omitempty"`
	Username                 string             `bson:"username"`
	Password                 string             `bson:"password"`
	Role                     Role               `bson:"role"`
	Locked                   bool               `bson:"locked"`
	LockReason               string             `bson:"lockReason"`
	Email                    string             `bson:"email"`
	Verified                 bool               `bson:"verified"`
	VerificationCode         string             `bson:"verificationCode"`
	VerificationCodeTime     time.Time          `bson:"verificationCodeTime"`
	NumOfErrTryLogin         int32              `bson:"numOfErrTryLogin"`
	LastErrTryLoginTime      time.Time          `bson:"lastErrTryLoginTime"`
	RecoveryPasswordCode     string             `bson:"recoveryPasswordCode"`
	RecoveryPasswordCodeTime time.Time          `bson:"recoveryPasswordCodeTime"`
}

func ConvertRoleToString(role Role) string {
	if role == 0 {
		return "USER"
	} else {
		return "ADMIN"
	}
}
