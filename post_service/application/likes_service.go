package application

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeService struct {
	store                    domain.PostStore
	authServiceAddress       string
	connectionServiceAddress string
	profileServiceAddress    string
}

func NewLikeService(postService *PostService) *LikeService {
	return &LikeService{
		store:                    postService.store,
		authServiceAddress:       postService.authServiceAddress,
		connectionServiceAddress: postService.connectionServiceAddress,
		profileServiceAddress:    postService.profileServiceAddress,
	}
}

func (service *LikeService) GiveLike(postId primitive.ObjectID, like *domain.Like) *domain.LikeDetailsDTO {
	err := service.store.InsertLike(postId, like)
	if err != nil {
		panic(fmt.Errorf("Invalid like"))
	}
	return service.getLikeDetails(postId, like)
}

func (service *LikeService) GetLike(postId primitive.ObjectID, likeId primitive.ObjectID) *domain.LikeDetailsDTO {
	like, err := service.store.GetLike(postId, likeId)
	if err != nil {
		panic(fmt.Errorf("Invalid like"))
	}
	return service.getLikeDetails(postId, like)
}

func (service *LikeService) GetLikesForPost(postId primitive.ObjectID) []*domain.LikeDetailsDTO {
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

func (service *LikeService) UndoLike(postId primitive.ObjectID, reactionId primitive.ObjectID) {
	service.store.UndoLike(postId, reactionId)
}

func (service *LikeService) getLikeDetails(postId primitive.ObjectID, like *domain.Like) *domain.LikeDetailsDTO {
	return &domain.LikeDetailsDTO{
		//Owner:       mapProfileToOwner(service.getPostOwnerProfile(ctx, like.OwnerId)),
		Like:   like,
		PostId: postId,
	}
}
