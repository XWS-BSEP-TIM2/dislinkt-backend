package application

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	nsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/notification_service_adapter"
	psa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/profile_service_adapter"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/util"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService struct {
	store                      domain.PostStore
	authServiceAdapter         asa.IAuthServiceAdapter
	connectionServiceAdapter   csa.IConnectionServiceAdapter
	profileServiceAdapter      psa.IProfileServiceAdapter
	notificationServiceAdapter nsa.INotificationServiceAdapter
	postAccessValidator        *util.PostAccessValidator
	ownerFinder                *util.OwnerFinder
}

func NewCommentService(postService *PostService) *CommentService {
	return &CommentService{
		store:                      postService.store,
		authServiceAdapter:         postService.authServiceAdapter,
		connectionServiceAdapter:   postService.connectionServiceAdapter,
		notificationServiceAdapter: postService.notificationServiceAdapter,
		profileServiceAdapter:      postService.profileServiceAdapter,
		postAccessValidator:        postService.postAccessValidator,
		ownerFinder:                postService.ownerFinder,
	}
}

func (service *CommentService) CreateComment(ctx context.Context, postId primitive.ObjectID, comment *domain.Comment) *domain.CommentDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx, comment.OwnerId)
	err := service.store.InsertComment(postId, comment)
	if err != nil {
		panic(fmt.Errorf("Invalid comment"))
	}

	commenter := service.profileServiceAdapter.GetSingleProfile(ctx, comment.OwnerId)
	post, _ := service.store.Get(postId)

	var notification pb.Notification
	notification.OwnerId = post.OwnerId.Hex()
	notification.ForwardUrl = "posts/" + postId.Hex()
	notification.Text = "commented your post"
	notification.UserFullName = commenter.Name + " " + commenter.Surname

	if notification.OwnerId != commenter.UserId.Hex() {
		service.notificationServiceAdapter.InsertNotification(ctx, &pb.InsertNotificationRequest{Notification: &notification})
	}

	return service.getCommentDetailsMapper(ctx, postId)(comment)
}

func (service *CommentService) GetComment(ctx context.Context, postId primitive.ObjectID, commentId primitive.ObjectID) *domain.CommentDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	comment, err := service.store.GetComment(postId, commentId)
	if err != nil {
		panic(fmt.Errorf("Invalid comment"))
	}
	return service.getCommentDetailsMapper(ctx, postId)(comment)
}

func (service *CommentService) GetCommentsForPost(ctx context.Context, postId primitive.ObjectID) []*domain.CommentDetailsDTO {
	service.postAccessValidator.ValidateUserAccessPost(ctx, postId)
	comments, err := service.store.GetCommentsForPost(postId)
	if err != nil {
		panic(fmt.Errorf("comments for post unavailable"))
	}
	commentsDetails, ok := funk.Map(comments, service.getCommentDetailsMapper(ctx, postId)).([]*domain.CommentDetailsDTO)
	if !ok {
		log("Error in conversion of comments to commentDetails")
		panic(fmt.Errorf("comments unavailable"))
	}
	return commentsDetails
}

func (service *CommentService) getCommentDetailsMapper(ctx context.Context, postId primitive.ObjectID) func(like *domain.Comment) *domain.CommentDetailsDTO {
	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx)
	return func(comment *domain.Comment) *domain.CommentDetailsDTO {
		return &domain.CommentDetailsDTO{
			Owner:   getOwner(comment.OwnerId),
			Comment: comment,
			PostId:  postId,
		}
	}
}
