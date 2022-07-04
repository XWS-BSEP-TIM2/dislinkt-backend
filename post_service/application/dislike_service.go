package application

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
	psa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/profile_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/util"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DislikeService struct {
	store                    domain.PostStore
	authServiceAdapter       asa.IAuthServiceAdapter
	connectionServiceAdapter csa.IConnectionServiceAdapter
	profileServiceAdapter    psa.IProfileServiceAdapter
	loggingServiceAdapter    lsa.ILoggingServiceAdapter
	postAccessValidator      *util.PostAccessValidator
	ownerFinder              *util.OwnerFinder
}

func NewDislikeService(postService *PostService) *DislikeService {
	return &DislikeService{
		store:                    postService.store,
		authServiceAdapter:       postService.authServiceAdapter,
		connectionServiceAdapter: postService.connectionServiceAdapter,
		profileServiceAdapter:    postService.profileServiceAdapter,
		loggingServiceAdapter:    postService.loggingServiceAdapter,
		postAccessValidator:      postService.postAccessValidator,
		ownerFinder:              postService.ownerFinder,
	}
}

func (service *DislikeService) GiveDislike(ctx context.Context, postId primitive.ObjectID, dislike *domain.Dislike) *domain.DislikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GiveDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx2, dislike.OwnerId)
	err := service.store.InsertDislike(postId, dislike)
	if err != nil {
		message := fmt.Sprintf("Error during dislike creation on post with id: %s", postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "GiveDislike", dislike.OwnerId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	message := fmt.Sprintf("User disliked post with id: %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GiveDislike", dislike.OwnerId.Hex(), message)

	return service.getDislikeDetailsMapper(ctx2, postId)(dislike)
}

func (service *DislikeService) GetDislike(ctx context.Context, postId primitive.ObjectID, dislikeId primitive.ObjectID) *domain.DislikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	dislike, err := service.store.GetDislike(postId, dislikeId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if err != nil {
		message := fmt.Sprintf("Disike with id: %s not found on post with id %s", dislikeId.Hex(), postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "WARNING", "GetDislike", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	message := fmt.Sprintf("User fetched dislike with id: %s on post with id %s", dislikeId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetDislike", requesterId.Hex(), message)
	return service.getDislikeDetailsMapper(ctx2, postId)(dislike)
}

func (service *DislikeService) GetDislikesForPost(ctx context.Context, postId primitive.ObjectID) []*domain.DislikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetDislikesForPost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	dislikes, err := service.store.GetDislikesForPost(postId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if err != nil {
		message := fmt.Sprintf("Dislikes on post with id %s unavailable.", postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "GetDislikesForPost", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	dislikeDetails, ok := funk.Map(dislikes, service.getDislikeDetailsMapper(ctx2, postId)).([]*domain.DislikeDetailsDTO)
	if !ok {
		log("Error in conversion of dislikes to commentDetails")
		panic(fmt.Errorf("dislikes unavailable"))
	}

	message := fmt.Sprintf("User fetched dislikes on post with id %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetDislikesForPost", requesterId.Hex(), message)
	return dislikeDetails
}

func (service *DislikeService) UndoDislike(ctx context.Context, postId primitive.ObjectID, reactionId primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "UndoDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	dislike, err := service.store.GetDislike(postId, reactionId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if err != nil {
		message := fmt.Sprintf("Cannot remove dislike with id: %s (not found)", reactionId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "UndoDislike", requesterId.Hex(), message)
		panic(errors.NewEntityNotFoundError(message))
	}
	service.authServiceAdapter.ValidateCurrentUser(ctx2, dislike.OwnerId)
	message := fmt.Sprintf("User removed dislike with id: %s from post with id: %s", reactionId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "UndoDislike", requesterId.Hex(), message)
	service.store.UndoDislike(postId, reactionId)
}

func (service *DislikeService) getDislikeDetailsMapper(ctx context.Context, postId primitive.ObjectID) func(dislike *domain.Dislike) *domain.DislikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "getDislikeDetailsMapper")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx2)
	return func(dislike *domain.Dislike) *domain.DislikeDetailsDTO {
		return &domain.DislikeDetailsDTO{
			Owner:   getOwner(dislike.OwnerId),
			Dislike: dislike,
			PostId:  postId,
		}
	}
}
