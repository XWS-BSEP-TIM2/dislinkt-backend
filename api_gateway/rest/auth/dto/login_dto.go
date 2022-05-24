package dto

type LoginRequestDto struct {
	Username string `json:"username" validate:"username_validation"`
	Password string `json:"password"`
}
