package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role int

const (
	USER Role = iota
	ADMIN
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Role     Role               `bson:"role"`
}

func ConvertRoleToString(role Role) string {
	if role == 0 {
		return "USER"
	} else {
		return "ADMIN"
	}
}
