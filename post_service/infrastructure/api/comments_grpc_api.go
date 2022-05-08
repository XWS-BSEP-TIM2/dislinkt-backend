package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentsSubHandler struct {
	service *application.CommentService
}

func NewCommentHandler(service *application.PostService) *CommentsSubHandler {
	return &CommentsSubHandler{
		service: application.NewCommentService(service),
	}
}

func (h *CommentsSubHandler) GetComment(ctx context.Context, request *pb.GetCommentRequest) (*pb.CommentResponse, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	commentId, err2 := primitive.ObjectIDFromHex(request.CommentId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given owner id is invalid."))
	}
	commentDetails := h.service.GetComment(postId, commentId)
	return &pb.CommentResponse{Comment: mapComment(commentDetails)}, nil

}

func (h *CommentsSubHandler) CreateComment(ctx context.Context, request *pb.CreateCommentRequest) (*pb.CommentResponse, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	ownerId, err2 := primitive.ObjectIDFromHex(request.NewComment.OwnerId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given comment id is invalid."))
	}
	commentDetails := h.service.CreateComment(postId, &domain.Comment{
		OwnerId:      ownerId,
		CreationTime: time.Now(),
		Content:      request.NewComment.Content,
	})
	return &pb.CommentResponse{Comment: mapComment(commentDetails)}, nil

}

func (h *CommentsSubHandler) GetComments(ctx context.Context, request *pb.GetPostRequest) (*pb.MultipleCommentsResponse, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	commentsDetails := h.service.GetCommentsForPost(postId)
	commentsResponse, ok := funk.Map(commentsDetails, func(dto *domain.CommentDetailsDTO) *pb.Comment { return mapComment(dto) }).([]*pb.Comment)
	if !ok {
		panic(fmt.Errorf("error during conversion of posts"))
	}
	return &pb.MultipleCommentsResponse{Comments: commentsResponse}, nil
}
