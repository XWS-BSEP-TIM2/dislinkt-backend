package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateProfileDto struct {
	FirstName string
	LastName  string
	Email     string
	BirthDate *timestamppb.Timestamp
	Gender    string
	IsPrivate bool
}
