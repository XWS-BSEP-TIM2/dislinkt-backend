package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostStore interface {
	Get(id primitive.ObjectID) (*Post, error)
	GetAll() ([]*Post, error)
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
}
