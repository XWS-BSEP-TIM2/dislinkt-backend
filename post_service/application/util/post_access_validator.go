package util

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostAccessValidator struct {
	store              domain.PostStore
	authServiceAdapter asa.IAuthServiceAdapter
	connServiceAdapter csa.IConnectionServiceAdapter
	loggServiceAdapter lsa.ILoggingServiceAdapter
}

func NewPostAccessValidator(
	store domain.PostStore,
	authAdapter asa.IAuthServiceAdapter,
	connAdapter csa.IConnectionServiceAdapter,
	loggAdapter lsa.ILoggingServiceAdapter) *PostAccessValidator {
	return &PostAccessValidator{
		store:              store,
		authServiceAdapter: authAdapter,
		connServiceAdapter: connAdapter,
		loggServiceAdapter: loggAdapter,
	}
}

func (validator *PostAccessValidator) ValidateUserAccessPost(ctx context.Context, postId primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "ValidateUserAccessPost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	post, err := validator.store.Get(postId)
	if err != nil {
		panic(errors.NewEntityNotFoundError("Post with given id does not exist."))
	}
	currentUserId := validator.authServiceAdapter.GetRequesterId(ctx2)

	if currentUserId != primitive.NilObjectID && currentUserId == post.OwnerId {
		return
	}

	res := validator.connServiceAdapter.CanUserAccessPostFromOwner(ctx2, currentUserId, post.OwnerId)
	if res {
		return
	}
	message := fmt.Sprintf("Current user (id: %s) is forbidden to access info from post with id: %s", currentUserId.Hex(), postId.Hex())
	validator.loggServiceAdapter.Log(ctx2, "WARNING", "ValidateUserAccessPost", currentUserId.Hex(), message)
	panic(errors.NewEntityForbiddenError(message))
}
