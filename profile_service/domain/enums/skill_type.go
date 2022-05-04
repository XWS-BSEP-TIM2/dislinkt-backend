package enums

import "strings"

type SkillType int

const (
	Skill SkillType = iota
	Interest
)

func (s SkillType) ToString() string {
	return [...]string{"Skill", "Interest"}[s]
}

func ToEnumSkill(skill string) SkillType {
	if strings.Compare(strings.ToLower(skill), "skill") == 0 {
		return Skill
	}
	return Interest
}
