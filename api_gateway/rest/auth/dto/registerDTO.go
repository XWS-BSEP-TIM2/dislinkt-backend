package dto

import (
	"fmt"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type RegisterDTO struct {
	Name            string `json:"name" validate:"name_validation"`
	Surname         string `json:"surname" validate:"name_validation"`
	Username        string `json:"username" validate:"username_validation"`
	Email           string `json:"email" validate:"email_validation"`
	Birthday        string `json:"birthday"`
	Gender          string `json:"gender"`
	PhoneNumber     string `json:"phoneNumber" validate:"number_validation"`
	IsPrivate       bool   `json:"isPrivate"`
	Password        string `json:"password" validate:"password_validation"`
	ConfirmPassword string `json:"confirmPassword" validate:"password_validation"`
}

func (r *RegisterDTO) ToProto(userID string) *pbProfile.Profile {

	t, err := time.Parse("2022-02-25", r.Birthday)
	if err != nil {
		fmt.Println("Error BirthDate")
	}

	return &pbProfile.Profile{
		Id:          userID,
		Name:        r.Name,
		Surname:     r.Surname,
		Username:    r.Username,
		Email:       r.Email,
		Biography:   "",
		Gender:      r.Gender,
		PhoneNumber: r.PhoneNumber,
		BirthDate:   timestamppb.New(t),
		IsPrivate:   r.IsPrivate,
		//Skills      []*Skill
		//Technologies []*Experience
		//Biography
	}
}
