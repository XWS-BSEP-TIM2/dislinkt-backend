package enums

type SkillType int

const (
	Skill SkillType = iota
	Interest
)

func (s SkillType) ToString() string {
	return [...]string{"Skill", "Interest"}[s]
}

func ToEnumSkill(skill string) SkillType {
	if skill == "Skill" {
		return Skill
	}
	return Interest
}
