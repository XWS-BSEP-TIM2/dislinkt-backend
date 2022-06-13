package api

import (
	converter "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/converter"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/domain/enums"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapProfile(profile *domain.Profile) *pb.Profile {
	profilePb := &pb.Profile{
		Id:          profile.Id.Hex(),
		Name:        profile.Name,
		Surname:     profile.Surname,
		Username:    profile.Username,
		Email:       profile.Email,
		Biography:   profile.Biography,
		Gender:      enums.Gender.ToString(profile.Gender),
		PhoneNumber: profile.PhoneNumber,
		BirthDate:   timestamppb.New(profile.BirthDate),
		IsPrivate:   profile.IsPrivate,
		IsTwoFactor: profile.IsTwoFactor,
	}

	for _, skill := range profile.Skills {
		profilePb.Skills = append(profilePb.Skills, &pb.Skill{
			Id:        skill.Id.Hex(),
			Name:      skill.Name,
			SkillType: enums.SkillType.ToString(skill.Type),
		})
	}

	for _, experience := range profile.Experiences {
		profilePb.Experiences = append(profilePb.Experiences, &pb.Experience{
			Id:             experience.Id.Hex(),
			Name:           experience.Name,
			Description:    experience.Description,
			ExperienceType: enums.ExperienceType.ToString(experience.TypeOfExperience),
			StartDate:      timestamppb.New(experience.StartDate),
			EndDate:        timestamppb.New(experience.EndDate),
		})
	}

	return profilePb
}

func MapProfile(profile *pb.Profile) domain.Profile {
	domainProfile := domain.Profile{
		Id:          converter.GetObjectId(profile.GetId()),
		Name:        profile.GetName(),
		Surname:     profile.GetSurname(),
		Username:    profile.GetUsername(),
		Email:       profile.GetEmail(),
		Biography:   profile.GetBiography(),
		Gender:      enums.ToEnumGender(profile.GetGender()),
		PhoneNumber: profile.PhoneNumber,
		BirthDate:   profile.GetBirthDate().AsTime(),
		IsPrivate:   profile.IsPrivate,
		Skills:      []domain.Skill{},
		Experiences: []domain.Experience{},
		IsTwoFactor: profile.IsTwoFactor,
	}

	for _, skill := range profile.Skills {
		domainProfile.Skills = append(domainProfile.Skills, domain.Skill{
			Id:   converter.GetObjectId(skill.Id),
			Name: skill.Name,
			Type: enums.ToEnumSkill(skill.SkillType),
		})
	}

	for _, experience := range profile.Experiences {
		domainProfile.Experiences = append(domainProfile.Experiences, domain.Experience{
			Id:               converter.GetObjectId(experience.Id),
			Name:             experience.Name,
			Description:      experience.Description,
			TypeOfExperience: enums.ToEnumExperience(experience.ExperienceType),
			StartDate:        experience.StartDate.AsTime(),
			EndDate:          experience.EndDate.AsTime(),
		})
	}

	return domainProfile
}
