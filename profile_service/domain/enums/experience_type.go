package enums

type ExperienceType int

const (
	Education ExperienceType = iota
	Working
)

func (e ExperienceType) ToString() string {
	return [...]string{"Education", "Working"}[e]
}

func ToEnumExperience(experience string) ExperienceType {
	if experience == "Education" {
		return Education
	}
	return Working
}
