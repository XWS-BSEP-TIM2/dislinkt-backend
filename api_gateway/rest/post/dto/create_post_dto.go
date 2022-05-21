package dto

type CreatePostDto struct {
	OwnerId     string   `json:"ownerId"`
	Content     string   `json:"content"`
	Links       []string `json:"links"`
	ImageBase64 string   `json:"imageBase64"`
}
