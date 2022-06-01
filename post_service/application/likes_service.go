package application

import (
	"context"
	"fmt"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
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
	profileServiceAddress    string
	postAccessValidator      *util.PostAccessValidator
}

func NewLikeService(postService *PostService) *LikeService {
	return &LikeService{
		store:                    postService.store,
		authServiceAdapter:       postService.authServiceAdapter,
		connectionServiceAdapter: postService.connectionServiceAdapter,
		profileServiceAddress:    postService.profileServiceAddress,
		postAccessValidator:      postService.postAccessValidator,
	}
}

func (service *LikeService) GiveLike(ctx context.Context, postId primitive.ObjectID, like *domain.Like) *domain.LikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx, like.OwnerId)
	err := service.store.InsertLike(postId, like)
	if err != nil {
		panic(fmt.Errorf("Invalid like"))
	}
	return service.getLikeDetails(postId, like)
}

func (service *LikeService) GetLike(ctx context.Context, postId primitive.ObjectID, likeId primitive.ObjectID) *domain.LikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	like, err := service.store.GetLike(postId, likeId)
	if err != nil {
		panic(fmt.Errorf("Invalid like"))
	}
	return service.getLikeDetails(postId, like)
}

func (service *LikeService) GetLikesForPost(ctx context.Context, postId primitive.ObjectID) []*domain.LikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	likes, err := service.store.GetLikesForPost(postId)
	if err != nil {
		panic(fmt.Errorf("likes for post unavailable"))
	}
	commentsDetails, ok := funk.Map(likes, func(like *domain.Like) *domain.LikeDetailsDTO {
		return &domain.LikeDetailsDTO{
			Like:   like,
			PostId: postId,
		}
	}).([]*domain.LikeDetailsDTO)
	if !ok {
		log("Error in conversion of likes to commentDetails")
		panic(fmt.Errorf("likes unavailable"))
	}
	return commentsDetails
}

func (service *LikeService) UndoLike(ctx context.Context, postId primitive.ObjectID, reactionId primitive.ObjectID) {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	like, err := service.store.GetLike(postId, reactionId)
	if err != nil {
		panic(errors.NewEntityNotFoundError("Cannot remove like with id: " + reactionId.Hex()))
	}
	service.authServiceAdapter.ValidateCurrentUser(ctx, like.OwnerId)
	service.store.UndoLike(postId, reactionId)
}

func (service *LikeService) getLikeDetails(postId primitive.ObjectID, like *domain.Like) *domain.LikeDetailsDTO {
	return &domain.LikeDetailsDTO{
		//Owner:       mapProfileToOwner(service.getPostOwnerProfile(ctx, like.OwnerId)),
		Like:   like,
		PostId: postId,
	}
}
