package connection_service_adapter

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/services"
	cb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionServiceAdapter struct {
	address string
}

func NewConnectionServiceAdapter(address string) *ConnectionServiceAdapter {
	return &ConnectionServiceAdapter{address: address}
}

func (conn *ConnectionServiceAdapter) GetAllUserConnections(ctx context.Context, id primitive.ObjectID) []*primitive.ObjectID {
	connClient := services.NewConnectionClient(conn.address)
	response, connErr := connClient.GetFriends(ctx, &cb.GetRequest{UserID: id.Hex()})
	if connErr != nil {
		panic(fmt.Errorf("Error during getting all connections: Connection Service"))
	}
	res, ok := funk.Map(response.Users, mapUserToUserId).([]*primitive.ObjectID)
	if !ok {
		panic(fmt.Errorf("Cannot cast list as []*primitive.ObjectId"))
	}
	return res
}

func (conn *ConnectionServiceAdapter) CanUserAccessPostFromOwner(ctx context.Context, userId primitive.ObjectID, ownerId primitive.ObjectID) bool {
	connClient := services.NewConnectionClient(conn.address)
	details, connErr := connClient.GetConnectionDetail(ctx, &cb.GetConnectionDetailRequest{
		UserIDa: userId.Hex(),
		UserIDb: ownerId.Hex(),
	})
	if connErr != nil {
		panic(fmt.Errorf("Error during getting connection details: Connection Service"))
	}
	if details.IsPrivate {
		if details.Relation != "CONNECTED" {
			return false
		}
	} else {
		if details.Relation == "B_BLOCK_A" {
			return false
		}
	}
	return true
}

func mapUserToUserId(user *cb.User) *primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(user.UserID)
	if err != nil {
		panic(fmt.Errorf("Error during id conversion while getting all connections: Connection Service"))
	}
	return &id
}
