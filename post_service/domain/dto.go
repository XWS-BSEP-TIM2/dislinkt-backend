package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostDetailsDTO struct {
	//Owner       *Owner
	Post        *Post
	ImageBase64 string
	Stats       *Stats
	Reactions   *Reactions
}

type Stats struct {
	CommentsNumber int
	LikesNumber    int
	DislikesNumber int
}

type Reactions struct {
	Liked    bool
	Disliked bool
}

type CommentDetailsDTO struct {
	//Owner       *Owner
	Comment *Comment
	PostId  primitive.ObjectID
}

type ReactionDetailsDTO struct {
	//Owner       *Owner
	ReactionId   primitive.ObjectID
	OwnerId      primitive.ObjectID
	CreationTime time.Time
	ReactionType string
	PostId       primitive.ObjectID
}

type LikeDetailsDTO struct {
	Like   *Like
	PostId primitive.ObjectID
}

type DislikeDetailsDTO struct {
	Dislike *Dislike
	PostId  primitive.ObjectID
}
