package application

import (
	"context"
	"fmt"
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
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx, like.OwnerId)
	err := service.store.InsertLike(postId, like)
	if err != nil {
		message := fmt.Sprintf("Error during like creation on post with id: %s", postId.Hex())
		service.loggingServiceAdapter.Log(ctx, "ERROR", "GiveLike", like.OwnerId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	message := fmt.Sprintf("User liked post with id: %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx, "SUCCESS", "GiveLike", like.OwnerId.Hex(), message)
	return service.getLikeDetailsMapper(ctx, postId)(like)
}

func (service *LikeService) GetLike(ctx context.Context, postId primitive.ObjectID, likeId primitive.ObjectID) *domain.LikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	like, err := service.store.GetLike(postId, likeId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx)
	if err != nil {
		message := fmt.Sprintf("Like with id: %s not found on post with id %s", likeId.Hex(), postId.Hex())
		service.loggingServiceAdapter.Log(ctx, "WARNING", "GetLike", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	message := fmt.Sprintf("User fetched like with id: %s on post with id %s", likeId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx, "SUCCESS", "GetLike", requesterId.Hex(), message)
	return service.getLikeDetailsMapper(ctx, postId)(like)
}

func (service *LikeService) GetLikesForPost(ctx context.Context, postId primitive.ObjectID) []*domain.LikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	likes, err := service.store.GetLikesForPost(postId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx)

	if err != nil {
		message := fmt.Sprintf("Likes on post with id %s unavailable.", postId.Hex())
		service.loggingServiceAdapter.Log(ctx, "ERROR", "GetLikesForPost", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	likeDetails, ok := funk.Map(likes, service.getLikeDetailsMapper(ctx, postId)).([]*domain.LikeDetailsDTO)
	if !ok {
		log("Error in conversion of likes to commentDetails")
		panic(fmt.Errorf("likes unavailable"))
	}
	message := fmt.Sprintf("User fetched likes on post with id %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx, "SUCCESS", "GetLikesForPost", requesterId.Hex(), message)
	return likeDetails
}

func (service *LikeService) UndoLike(ctx context.Context, postId primitive.ObjectID, reactionId primitive.ObjectID) {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	like, err := service.store.GetLike(postId, reactionId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx)

	if err != nil {
		message := fmt.Sprintf("Cannot remove like with id: %s (not found)", reactionId.Hex())
		service.loggingServiceAdapter.Log(ctx, "ERROR", "UndoLike", requesterId.Hex(), message)
		panic(errors.NewEntityNotFoundError(message))
	}
	service.authServiceAdapter.ValidateCurrentUser(ctx, like.OwnerId)
	message := fmt.Sprintf("User removed like with id: %s from post with id: %s", reactionId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx, "SUCCESS", "UndoLike", requesterId.Hex(), message)
	service.store.UndoLike(postId, reactionId)
}

func (service *LikeService) getLikeDetailsMapper(ctx context.Context, postId primitive.ObjectID) func(like *domain.Like) *domain.LikeDetailsDTO {
	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx)
	return func(like *domain.Like) *domain.LikeDetailsDTO {
		return &domain.LikeDetailsDTO{
			Owner:  getOwner(like.OwnerId),
			Like:   like,
			PostId: postId,
		}
	}
}
