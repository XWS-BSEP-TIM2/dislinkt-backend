package application

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	store                    domain.PostStore
	authServiceAddress       string
	connectionServiceAddress string
	profileServiceAddress    string
}

func NewCommentService(postService *PostService) *CommentService {
	return &CommentService{
		store:                    postService.store,
		authServiceAddress:       postService.authServiceAddress,
		connectionServiceAddress: postService.connectionServiceAddress,
		profileServiceAddress:    postService.profileServiceAddress,
	}
}

func (service *CommentService) CreateComment(postId primitive.ObjectID, comment *domain.Comment) *domain.CommentDetailsDTO {
	err := service.store.InsertComment(postId, comment)
	if err != nil {
		panic(fmt.Errorf("Invalid comment"))
	}
	return service.getCommentDetails(postId, comment)
}

func (service *CommentService) GetComment(postId primitive.ObjectID, commentId primitive.ObjectID) *domain.CommentDetailsDTO {
	comment, err := service.store.GetComment(postId, commentId)
	if err != nil {
		panic(fmt.Errorf("Invalid comment"))
	}
	return service.getCommentDetails(postId, comment)
}

func (service *CommentService) GetCommentsForPost(postId primitive.ObjectID) []*domain.CommentDetailsDTO {
	comments, err := service.store.GetCommentsForPost(postId)
	if err != nil {
		panic(fmt.Errorf("comments for post unavailable"))
	}
	commentsDetails, ok := funk.Map(comments, func(comment *domain.Comment) *domain.CommentDetailsDTO {
		return &domain.CommentDetailsDTO{
			Comment: comment,
			PostId:  postId,
		}
	}).([]*domain.CommentDetailsDTO)
	if !ok {
		log("Error in conversion of comments to commentDetails")
		panic(fmt.Errorf("comments unavailable"))
	}
	return commentsDetails
}

func (service *CommentService) getCommentDetails(postId primitive.ObjectID, comment *domain.Comment) *domain.CommentDetailsDTO {
	return &domain.CommentDetailsDTO{
		//Owner:       mapProfileToOwner(service.getPostOwnerProfile(ctx, comment.OwnerId)),
		Comment: comment,
		PostId:  postId,
	}
}
