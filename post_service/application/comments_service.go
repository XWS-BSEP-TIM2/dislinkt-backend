package application

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	asa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/auth_service_adapter"
	csa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/connection_service_adapter"
	lsa "github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application/adapters/logging_service_adapter"
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
	loggingServiceAdapter      lsa.ILoggingServiceAdapter
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
		loggingServiceAdapter:      postService.loggingServiceAdapter,
		postAccessValidator:        postService.postAccessValidator,
		ownerFinder:                postService.ownerFinder,
	}
}

func (service *CommentService) CreateComment(ctx context.Context, postId primitive.ObjectID, comment *domain.Comment) *domain.CommentDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "CreateComment")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	service.authServiceAdapter.ValidateCurrentUser(ctx2, comment.OwnerId)
	err := service.store.InsertComment(postId, comment)
	if err != nil {
		message := fmt.Sprintf("Error during comment creation on post with id: %s", postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "CreateComment", comment.OwnerId.Hex(), message)
		panic(fmt.Errorf(message))
	}

	commenter := service.profileServiceAdapter.GetSingleProfile(ctx2, comment.OwnerId)
	post, _ := service.store.Get(postId)

	var notification pb.Notification
	notification.OwnerId = post.OwnerId.Hex()
	notification.ForwardUrl = "posts/" + postId.Hex()
	notification.Text = "commented on your post"
	notification.UserFullName = commenter.Name + " " + commenter.Surname

	if notification.OwnerId != commenter.UserId.Hex() {
		service.notificationServiceAdapter.InsertNotification(ctx2, &pb.InsertNotificationRequest{Notification: &notification})
	}

	message := fmt.Sprintf("User commented on post with id: %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "CreateComment", comment.OwnerId.Hex(), message)
	return service.getCommentDetailsMapper(ctx2, postId)(comment)
}

func (service *CommentService) GetComment(ctx context.Context, postId primitive.ObjectID, commentId primitive.ObjectID) *domain.CommentDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetComment")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	comment, err := service.store.GetComment(postId, commentId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)

	if err != nil {
		message := fmt.Sprintf("Comment with id: %s not found on post with id %s", commentId.Hex(), postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "WARNING", "GetComment", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}

	message := fmt.Sprintf("User fetched comment with id: %s on post with id %s", commentId.Hex(), postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetComment", requesterId.Hex(), message)
	return service.getCommentDetailsMapper(ctx2, postId)(comment)
}

func (service *CommentService) GetCommentsForPost(ctx context.Context, postId primitive.ObjectID) []*domain.CommentDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "GetCommentsForPost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	service.postAccessValidator.ValidateUserAccessPost(ctx2, postId)
	comments, err := service.store.GetCommentsForPost(postId)
	requesterId := service.authServiceAdapter.GetRequesterId(ctx2)
	if err != nil {
		message := fmt.Sprintf("Comments on post with id %s unavailable.", postId.Hex())
		service.loggingServiceAdapter.Log(ctx2, "ERROR", "GetCommentsForPost", requesterId.Hex(), message)
		panic(fmt.Errorf(message))
	}
	commentsDetails, ok := funk.Map(comments, service.getCommentDetailsMapper(ctx2, postId)).([]*domain.CommentDetailsDTO)
	if !ok {
		panic(fmt.Errorf("Error during mapping of CommentDetails"))
	}

	message := fmt.Sprintf("User fetched comments on post with id %s", postId.Hex())
	service.loggingServiceAdapter.Log(ctx2, "SUCCESS", "GetCommentsForPost", requesterId.Hex(), message)
	return commentsDetails
}

func (service *CommentService) getCommentDetailsMapper(ctx context.Context, postId primitive.ObjectID) func(like *domain.Comment) *domain.CommentDetailsDTO {
	span := tracer.StartSpanFromContext(ctx, "getCommentDetailsMapper")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	getOwner := service.ownerFinder.GetOwnerFinderFunction(ctx2)
	return func(comment *domain.Comment) *domain.CommentDetailsDTO {
		return &domain.CommentDetailsDTO{
			Owner:   getOwner(comment.OwnerId),
			Comment: comment,
			PostId:  postId,
		}
	}
}
