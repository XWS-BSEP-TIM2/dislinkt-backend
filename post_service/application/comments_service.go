package application

import (
	"context"
	"fmt"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/util"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	store                    domain.PostStore
	authServiceAdapter       asa.IAuthServiceAdapter
	connectionServiceAdapter csa.IConnectionServiceAdapter
	profileServiceAddress    string
	postAccessValidator      *util.PostAccessValidator
}

func NewCommentService(postService *PostService) *CommentService {
	return &CommentService{
		store:                    postService.store,
		authServiceAdapter:       postService.authServiceAdapter,
		connectionServiceAdapter: postService.connectionServiceAdapter,
		profileServiceAddress:    postService.profileServiceAddress,
		postAccessValidator:      postService.postAccessValidator,
	}
}

func (service *CommentService) CreateComment(ctx context.Context, postId primitive.ObjectID, comment *domain.Comment) *domain.CommentDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx, comment.OwnerId)
	err := service.store.InsertComment(postId, comment)
	if err != nil {
		panic(fmt.Errorf("Invalid comment"))
	}
	return service.getCommentDetails(postId, comment)
}

func (service *CommentService) GetComment(ctx context.Context, postId primitive.ObjectID, commentId primitive.ObjectID) *domain.CommentDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	comment, err := service.store.GetComment(postId, commentId)
	if err != nil {
		panic(fmt.Errorf("Invalid comment"))
	}
	return service.getCommentDetails(postId, comment)
}

func (service *CommentService) GetCommentsForPost(ctx context.Context, postId primitive.ObjectID) []*domain.CommentDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
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
