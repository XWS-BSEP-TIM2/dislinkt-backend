package DTO

import (
	"fmt"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type RegisterDTO struct {
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Birthday        string `json:"birthday"`
	Gender          string `json:"gender"`
	PhoneNumber     string `json:"phoneNumber"`
	IsPrivate       bool   `json:"isPrivate"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
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
		//Experiences []*Experience
		//Biography
	}
}
