package dto

type CreateCommentDto struct {
	OwnerId string `json:"ownerId"`
	Content string `json:"content"`
}
