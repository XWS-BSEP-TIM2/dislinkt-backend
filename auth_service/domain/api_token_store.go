package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApiTokenStore interface {
	Get(ctx context.Context, id primitive.ObjectID) (*ApiToken, error)
	Insert(ctx context.Context, product *ApiToken) (error, string)
	DeleteById(ctx context.Context, id primitive.ObjectID) (int64, error)
	GetByTokenCode(ctx context.Context, tokenCode string) (*ApiToken, error)
	DeleteAllUserTokens(ctx context.Context, id primitive.ObjectID) error
	DeleteAllTokens(ctx context.Context) error
}
