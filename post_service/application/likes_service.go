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

type LikeService struct {
	store                    domain.PostStore
	authServiceAdapter       asa.IAuthServiceAdapter
	connectionServiceAdapter csa.IConnectionServiceAdapter
	profileServiceAdapter    psa.IProfileServiceAdapter
	loggingServiceAdapter    lsa.ILoggingServiceAdapter
	postAccessValidator      *util.PostAccessValidator
	ownerFinder              *util.OwnerFinder
}

func NewLikeService(postService *PostService) *LikeService {
	return &LikeService{
		store:                    postService.store,
		authServiceAdapter:       postService.authServiceAdapter,
		connectionServiceAdapter: postService.connectionServiceAdapter,
		profileServiceAdapter:    postService.profileServiceAdapter,
		loggingServiceAdapter:    postService.loggingServiceAdapter,
		postAccessValidator:      postService.postAccessValidator,
		ownerFinder:              postService.ownerFinder,
	}
}

func (service *LikeService) GiveLike(ctx context.Context, postId primitive.ObjectID, like *domain.Like) *domain.LikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GiveLike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx2, like.OwnerId)
	err := service.store.InsertLike(postId, like)
	if err != nil {
		message := fmt.Sprintf("Error during like creation on post with id: %s", postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "GiveLike", like.OwnerId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	message := fmt.Sprintf("User liked post with id: %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GiveLike", like.OwnerId.Hex(), message)
	return service.getLikeDetailsMapper(ctx2, postId)(like)
}

func (service *LikeService) GetLike(ctx context.Context, postId primitive.ObjectID, likeId primitive.ObjectID) *domain.LikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetLike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	like, err := service.store.GetLike(postId, likeId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)
	if err != nil {
		message := fmt.Sprintf("Like with id: %s not found on post with id %s", likeId.Hex(), postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "WARNING", "GetLike", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	message := fmt.Sprintf("User fetched like with id: %s on post with id %s", likeId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetLike", requesterId.Hex(), message)
	return service.getLikeDetailsMapper(ctx2, postId)(like)
}

func (service *LikeService) GetLikesForPost(ctx context.Context, postId primitive.ObjectID) []*domain.LikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetLikesForPost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	likes, err := service.store.GetLikesForPost(postId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if err != nil {
		message := fmt.Sprintf("Likes on post with id %s unavailable.", postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "GetLikesForPost", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	likeDetails, ok := funk.Map(likes, service.getLikeDetailsMapper(ctx2, postId)).([]*domain.LikeDetailsDTO)
	if !ok {
		log("Error in conversion of likes to commentDetails")
		panic(fmt.Errorf("likes unavailable"))
	}
	message := fmt.Sprintf("User fetched likes on post with id %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetLikesForPost", requesterId.Hex(), message)
	return likeDetails
}

func (service *LikeService) UndoLike(ctx context.Context, postId primitive.ObjectID, reactionId primitive.ObjectID) {
	span := tracer.StartSpanFromContext(ctx, "UndoLike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	like, err := service.store.GetLike(postId, reactionId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if err != nil {
		message := fmt.Sprintf("Cannot remove like with id: %s (not found)", reactionId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "UndoLike", requesterId.Hex(), message)
		panic(errors.NewEntityNotFoundError(message))
	}
	service.authServiceAdapter.ValidateCurrentUser(ctx2, like.OwnerId)
	message := fmt.Sprintf("User removed like with id: %s from post with id: %s", reactionId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "UndoLike", requesterId.Hex(), message)
	service.store.UndoLike(postId, reactionId)
}

func (service *LikeService) getLikeDetailsMapper(ctx context.Context, postId primitive.ObjectID) func(like *domain.Like) *domain.LikeDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "getLikeDetailsMapper")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx2)
	return func(like *domain.Like) *domain.LikeDetailsDTO {
		return &domain.LikeDetailsDTO{
			Owner:  getOwner(like.OwnerId),
			Like:   like,
			PostId: postId,
		}
	}
}
