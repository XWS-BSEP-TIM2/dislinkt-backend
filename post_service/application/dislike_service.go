package application

import (
	"context"
	"fmt"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
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
	postAccessValidator      *util.PostAccessValidator
	ownerFinder              *util.OwnerFinder
}

func NewDislikeService(postService *PostService) *DislikeService {
	return &DislikeService{
		store:                    postService.store,
		authServiceAdapter:       postService.authServiceAdapter,
		connectionServiceAdapter: postService.connectionServiceAdapter,
		profileServiceAdapter:    postService.profileServiceAdapter,
		postAccessValidator:      postService.postAccessValidator,
		ownerFinder:              postService.ownerFinder,
	}
}

func (service *DislikeService) GiveDislike(ctx context.Context, postId primitive.ObjectID, dislike *domain.Dislike) *domain.DislikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx, dislike.OwnerId)
	err := service.store.InsertDislike(postId, dislike)
	if err != nil {
		panic(fmt.Errorf("Invalid dislike"))
	}
	return service.getDislikeDetailsMapper(ctx, postId)(dislike)
}

func (service *DislikeService) GetDislike(ctx context.Context, postId primitive.ObjectID, dislikeId primitive.ObjectID) *domain.DislikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	dislike, err := service.store.GetDislike(postId, dislikeId)
	if err != nil {
		panic(fmt.Errorf("Invalid dislike"))
	}
	return service.getDislikeDetailsMapper(ctx, postId)(dislike)
}

func (service *DislikeService) GetDislikesForPost(ctx context.Context, postId primitive.ObjectID) []*domain.DislikeDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	dislikes, err := service.store.GetDislikesForPost(postId)
	if err != nil {
		panic(fmt.Errorf("dislikes for post unavailable"))
	}
	dislikeDetails, ok := funk.Map(dislikes, service.getDislikeDetailsMapper(ctx, postId)).([]*domain.DislikeDetailsDTO)
	if !ok {
		log("Error in conversion of dislikes to commentDetails")
		panic(fmt.Errorf("dislikes unavailable"))
	}
	return dislikeDetails
}

func (service *DislikeService) UndoDislike(ctx context.Context, postId primitive.ObjectID, reactionId primitive.ObjectID) {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	dislike, err := service.store.GetDislike(postId, reactionId)
	if err != nil {
		panic(errors.NewEntityNotFoundError("Cannot remove dislike with id: " + reactionId.Hex()))
	}
	service.authServiceAdapter.ValidateCurrentUser(ctx, dislike.OwnerId)
	service.store.UndoDislike(postId, reactionId)
}

func (service *DislikeService) getDislikeDetailsMapper(ctx context.Context, postId primitive.ObjectID) func(dislike *domain.Dislike) *domain.DislikeDetailsDTO {
	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx)
	return func(dislike *domain.Dislike) *domain.DislikeDetailsDTO {
		return &domain.DislikeDetailsDTO{
			Owner:   getOwner(dislike.OwnerId),
			Dislike: dislike,
			PostId:  postId,
		}
	}
}
