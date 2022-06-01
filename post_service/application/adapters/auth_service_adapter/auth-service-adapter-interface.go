package auth_service_adapter

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthServiceAdapter interface {
	GetRequesterId(ctx context.Context) primitive.ObjectID
	ValidateCurrentUser(ctx context.Context, userId primitive.ObjectID)
}
