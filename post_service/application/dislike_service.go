package application

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DislikeService struct {
	store                    domain.PostStore
	authServiceAddress       string
	connectionServiceAddress string
	profileServiceAddress    string
}

func NewDislikeService(postService *PostService) *DislikeService {
	return &DislikeService{
		store:                    postService.store,
		authServiceAddress:       postService.authServiceAddress,
		connectionServiceAddress: postService.connectionServiceAddress,
		profileServiceAddress:    postService.profileServiceAddress,
	}
}

func (service *DislikeService) GiveDislike(postId primitive.ObjectID, dislike *domain.Dislike) *domain.DislikeDetailsDTO {
	err := service.store.InsertDislike(postId, dislike)
	if err != nil {
		panic(fmt.Errorf("Invalid dislike"))
	}
	return service.getDislikeDetails(postId, dislike)
}

func (service *DislikeService) GetDislike(postId primitive.ObjectID, dislikeId primitive.ObjectID) *domain.DislikeDetailsDTO {
	dislike, err := service.store.GetDislike(postId, dislikeId)
	if err != nil {
		panic(fmt.Errorf("Invalid dislike"))
	}
	return service.getDislikeDetails(postId, dislike)
}

func (service *DislikeService) GetDislikesForPost(postId primitive.ObjectID) []*domain.DislikeDetailsDTO {
	dislikes, err := service.store.GetDislikesForPost(postId)
	if err != nil {
		panic(fmt.Errorf("dislikes for post unavailable"))
	}
	dislikeDetails, ok := funk.Map(dislikes, func(dislike *domain.Dislike) *domain.DislikeDetailsDTO {
		return &domain.DislikeDetailsDTO{
			Dislike: dislike,
			PostId:  postId,
		}
	}).([]*domain.DislikeDetailsDTO)
	if !ok {
		log("Error in conversion of dislikes to commentDetails")
		panic(fmt.Errorf("dislikes unavailable"))
	}
	return dislikeDetails
}

func (service *DislikeService) UndoDislike(postId primitive.ObjectID, reactionId primitive.ObjectID) {
	service.store.UndoDislike(postId, reactionId)
}

func (service *DislikeService) getDislikeDetails(postId primitive.ObjectID, dislike *domain.Dislike) *domain.DislikeDetailsDTO {
	return &domain.DislikeDetailsDTO{
		//Owner:       mapProfileToOwner(service.getPostOwnerProfile(ctx, dislike.OwnerId)),
		Dislike: dislike,
		PostId:  postId,
	}
}
