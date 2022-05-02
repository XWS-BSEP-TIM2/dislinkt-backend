package enums

type Gender int

func ToEnumGender(gender string) Gender {
	if gender == "Male" {
		return Male
	} else {
		return Female
	}
}

const (
	Male Gender = iota
	Female
)

func (g Gender) ToString() string {
	return [...]string{"Male", "Female"}[g]
}
