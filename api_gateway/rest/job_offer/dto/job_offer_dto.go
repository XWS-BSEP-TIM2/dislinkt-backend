package dto

type JobOfferDto struct {
	Id           string   `json:"id"`
	CompanyName  string   `json:"companyName"`
	UserId       string   `json:"userId"`
	Position     string   `json:"position"`
	Seniority    string   `json:"seniority"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
}
