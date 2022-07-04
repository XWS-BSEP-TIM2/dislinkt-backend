package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CommentSubHandler struct {
	service *application.CommentService
}

func NewCommentHandler(service *application.PostService) *CommentSubHandler {
	return &CommentSubHandler{
		service: application.NewCommentService(service),
	}
}

func (h *CommentSubHandler) GetComment(ctx context.Context, request *pb.GetSubresourceRequest) (*pb.CommentResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetComment")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	commentId, err2 := primitive.ObjectIDFromHex(request.SubresourceId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given owner id is invalid."))
	}
	commentDetails := h.service.GetComment(ctx2, postId, commentId)
	return &pb.CommentResponse{Comment: mapComment(commentDetails)}, nil

}

func (h *CommentSubHandler) CreateComment(ctx context.Context, request *pb.CreateCommentRequest) (*pb.CommentResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateComment")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	ownerId, err2 := primitive.ObjectIDFromHex(request.NewComment.OwnerId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given owner id is invalid."))
	}
	commentDetails := h.service.CreateComment(ctx2, postId, &domain.Comment{
		OwnerId:      ownerId,
		CreationTime: time.Now(),
		Content:      request.NewComment.Content,
	})
	return &pb.CommentResponse{Comment: mapComment(commentDetails)}, nil

}

func (h *CommentSubHandler) GetComments(ctx context.Context, request *pb.GetPostRequest) (*pb.MultipleCommentsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetComments")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	commentsDetails := h.service.GetCommentsForPost(ctx2, postId)
	commentsResponse, ok := funk.Map(commentsDetails, func(dto *domain.CommentDetailsDTO) *pb.Comment { return mapComment(dto) }).([]*pb.Comment)
	if !ok {
		panic(fmt.Errorf("error during conversion of posts"))
	}
	return &pb.MultipleCommentsResponse{Comments: commentsResponse}, nil
}
