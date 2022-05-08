package DTO

import (
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
)

type ConnectionDTO struct {
	UserID    string
	Name      string
	Surname   string
	Username  string
	Biography string
	IsPrivate bool
	Skills    []Skill
}

type Skill struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"skillType"`
}

func MapConnectionDTO(profile *pbProfile.Profile) *ConnectionDTO {

	return &ConnectionDTO{
		UserID:    profile.Id,
		Name:      profile.Name,
		Surname:   profile.Surname,
		Username:  profile.Username,
		Biography: profile.Biography,
		IsPrivate: profile.IsPrivate,
		Skills:    GetSkills(profile.Skills),
	}
}

func GetSkills(pbSkills []*pbProfile.Skill) []Skill {
	var skills []Skill

	for _, v := range pbSkills {
		skills = append(skills, Skill{Id: v.Id, Name: v.Name, Type: v.SkillType})
	}

	return skills
}
