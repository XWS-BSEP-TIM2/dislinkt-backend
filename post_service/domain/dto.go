package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostDetailsDTO struct {
	//Owner       *Owner
	Post        *Post
	ImageBase64 string
	Stats       *Stats
}

type Stats struct {
	CommentsNumber int
	LikesNumber    int
	DislikesNumber int
}

type CommentDetailsDTO struct {
	//Owner       *Owner
	Comment *Comment
	PostId  primitive.ObjectID
}
