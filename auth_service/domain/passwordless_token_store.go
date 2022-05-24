package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordlessTokenStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*PasswordlessToken, error)
	Insert(ctx context.Context, product *PasswordlessToken) (error, string)
	DeleteById(ctx context.Context, id primitive.ObjectID) (int64, error)
	GetByTokenCode(ctx context.Context, tokenCode string) (*PasswordlessToken, error)
}
