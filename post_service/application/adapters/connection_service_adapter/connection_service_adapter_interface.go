package connection_service_adapter

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IConnectionServiceAdapter interface {
	GetAllUserConnections(ctx context.Context, id primitive.ObjectID) []*primitive.ObjectID
	CanUserAccessPostFromOwner(ctx context.Context, userId primitive.ObjectID, ownerId primitive.ObjectID) bool
}
