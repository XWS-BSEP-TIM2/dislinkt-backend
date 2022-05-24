package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Insert(ctx context.Context, product *User) (error, string)
	DeleteAll(ctx context.Context)
	GetByUsername(ctx context.Context, Username string) (*User, error)
	Update(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, Email string) (*User, error)
}
