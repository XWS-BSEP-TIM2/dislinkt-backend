package domain

type JobOffer struct {
	Id                 string
	UserId             string
	JobOfferUniqueCode string
	CompanyName        string
	Position           string
	Seniority          string
	Description        string
	Technologies       []string
}
