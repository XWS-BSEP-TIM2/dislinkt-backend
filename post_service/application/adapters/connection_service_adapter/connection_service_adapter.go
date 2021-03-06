package connection_service_adapter

import (
	"context"
	"fmt"
	cb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectionServiceAdapter struct {
	address               string
	loggingServiceAdapter lsa.ILoggingServiceAdapter
}

func NewConnectionServiceAdapter(address string, loggingServiceAdapter lsa.ILoggingServiceAdapter) *ConnectionServiceAdapter {
	return &ConnectionServiceAdapter{address: address, loggingServiceAdapter: loggingServiceAdapter}
}

func (conn *ConnectionServiceAdapter) GetAllUserConnections(ctx context.Context, id primitive.ObjectID) []*primitive.ObjectID {
	span := tracer.StartSpanFromContext(ctx, "GetAllUserConnections")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	connClient := adapters.NewConnectionClient(conn.address)
	response, connErr := connClient.GetFriends(ctx2, &cb.GetRequest{UserID: id.Hex()})
	if connErr != nil {
		message := "Error during getting all connections: Connection Service"
		conn.loggingServiceAdapter.Log(ctx2, "ERROR", "GetAllUserConnections", id.Hex(), message)
		panic(fmt.Errorf(message))
	}
	res, ok := funk.Map(response.Users, mapUserToUserId).([]*primitive.ObjectID)
	if !ok {
		panic(fmt.Errorf("Cannot cast list as []*primitive.ObjectId"))
	}
	return res
}

func (conn *ConnectionServiceAdapter) CanUserAccessPostFromOwner(ctx context.Context, userId primitive.ObjectID, ownerId primitive.ObjectID) bool {
	span := tracer.StartSpanFromContext(ctx, "CanUserAccessPostFromOwner")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	connClient := adapters.NewConnectionClient(conn.address)
	details, connErr := connClient.GetConnectionDetail(ctx2, &cb.GetConnectionDetailRequest{
		UserIDa: userId.Hex(),
		UserIDb: ownerId.Hex(),
	})
	if connErr != nil {
		message := "Error during getting connection details: Connection Service"
		conn.loggingServiceAdapter.Log(ctx2, "ERROR", "CanUserAccessPostFromOwner", userId.Hex(), message)
		panic(fmt.Errorf(message))
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
