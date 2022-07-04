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

type DislikeSubHandler struct {
	service *application.DislikeService
}

func NewDislikeHandler(service *application.PostService) *DislikeSubHandler {
	return &DislikeSubHandler{
		service: application.NewDislikeService(service),
	}
}

func (h *DislikeSubHandler) GetDislike(ctx context.Context, request *pb.GetSubresourceRequest) (*pb.ReactionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	reactionId, err2 := primitive.ObjectIDFromHex(request.SubresourceId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given reaction id is invalid."))
	}
	likeDetails := h.service.GetDislike(ctx2, postId, reactionId)
	return &pb.ReactionResponse{Reaction: mapDislike(likeDetails)}, nil

}

func (h *DislikeSubHandler) GiveDislike(ctx context.Context, request *pb.CreateReactionRequest) (*pb.ReactionResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GiveDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	ownerId, err2 := primitive.ObjectIDFromHex(request.NewReaction.OwnerId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given owner id is invalid."))
	}
	dislikeDetails := h.service.GiveDislike(ctx2, postId, &domain.Dislike{
		OwnerId:      ownerId,
		CreationTime: time.Now(),
	})
	return &pb.ReactionResponse{Reaction: mapDislike(dislikeDetails)}, nil

}

func (h *DislikeSubHandler) GetDislikes(ctx context.Context, request *pb.GetPostRequest) (*pb.MultipleReactionsResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetDislikes")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	dislikes := h.service.GetDislikesForPost(ctx2, postId)
	reactionsResponse, ok := funk.Map(dislikes, func(dto *domain.DislikeDetailsDTO) *pb.Reaction { return mapDislike(dto) }).([]*pb.Reaction)
	if !ok {
		panic(fmt.Errorf("error during conversion of posts"))
	}
	return &pb.MultipleReactionsResponse{Reactions: reactionsResponse}, nil
}

func (h *DislikeSubHandler) UndoDislike(ctx context.Context, request *pb.GetSubresourceRequest) (*pb.EmptyRequest, error) {
	span := tracer.StartSpanFromContext(ctx, "UndoDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	reactionId, err2 := primitive.ObjectIDFromHex(request.SubresourceId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given reaction id is invalid."))
	}
	h.service.UndoDislike(ctx2, postId, reactionId)
	return &pb.EmptyRequest{}, nil
}
