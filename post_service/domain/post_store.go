package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostStore interface {
	Get(id primitive.ObjectID) (*Post, error)
	GetAll() ([]*Post, error)
	GetPostsFromUser(userId primitive.ObjectID) ([]*Post, error)
	Insert(post *Post) error
	Update(post *Post) error
	InsertComment(postId primitive.ObjectID, comment *Comment) error
	GetComment(postId primitive.ObjectID, commentId primitive.ObjectID) (*Comment, error)
	DeleteAll()
	GetCommentsForPost(postId primitive.ObjectID) ([]*Comment, error)
	InsertLike(postId primitive.ObjectID, like *Like) error
	GetLike(postId primitive.ObjectID, likeId primitive.ObjectID) (*Like, error)
	GetLikesForPost(postId primitive.ObjectID) ([]*Like, error)
	UndoLike(postId primitive.ObjectID, reactionId primitive.ObjectID) error
	InsertDislike(postId primitive.ObjectID, dislike *Dislike) error
	GetDislike(postId primitive.ObjectID, dislikeId primitive.ObjectID) (*Dislike, error)
	GetDislikesForPost(postId primitive.ObjectID) ([]*Dislike, error)
	UndoDislike(postId primitive.ObjectID, reactionId primitive.ObjectID) error
}
