package util

import (
	"context"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostAccessValidator struct {
	store              domain.PostStore
	authServiceAdapter *asa.AuthServiceAdapter
	connServiceAddress *csa.ConnectionServiceAdapter
}

func NewPostAccessValidator(store domain.PostStore, authAdapter *asa.AuthServiceAdapter, connAdapter *csa.ConnectionServiceAdapter) *PostAccessValidator {
	return &PostAccessValidator{
		store:              store,
		authServiceAdapter: authAdapter,
		connServiceAddress: connAdapter,
	}
}

func (validator *PostAccessValidator) ValidateUserAccessPost(ctx context.Context, postId primitive.ObjectID) {
	post, err := validator.store.Get(postId)
	if err != nil {
		panic(errors.NewEntityNotFoundError("Post with given id does not exist."))
	}
	currentUserId := validator.authServiceAdapter.GetRequesterId(ctx)

	if currentUserId != primitive.NilObjectID && currentUserId == post.OwnerId {
		return
	}

	res := validator.connServiceAddress.CanUserAccessPostFromOwner(ctx, currentUserId, post.OwnerId)
	if res {
		return
	}
	panic(errors.NewEntityForbiddenError("Current user cannot access info from post with id: " + postId.Hex()))
}
