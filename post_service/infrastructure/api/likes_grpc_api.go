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

type LikeSubHandler struct {
	service *application.LikeService
}

func NewLikeHandler(service *application.PostService) *LikeSubHandler {
	return &LikeSubHandler{
		service: application.NewLikeService(service),
	}
}

func (h *LikeSubHandler) GetLike(ctx context.Context, request *pb.GetSubresourceRequest) (*pb.ReactionResponse, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	reactionId, err2 := primitive.ObjectIDFromHex(request.SubresourceId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given reaction id is invalid."))
	}
	likeDetails := h.service.GetLike(ctx, postId, reactionId)
	return &pb.ReactionResponse{Reaction: mapLike(likeDetails)}, nil

}

func (h *LikeSubHandler) GiveLike(ctx context.Context, request *pb.CreateReactionRequest) (*pb.ReactionResponse, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	ownerId, err2 := primitive.ObjectIDFromHex(request.NewReaction.OwnerId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given owner id is invalid."))
	}
	likeDetails := h.service.GiveLike(ctx, postId, &domain.Like{
		OwnerId:      ownerId,
		CreationTime: time.Now(),
	})
	return &pb.ReactionResponse{Reaction: mapLike(likeDetails)}, nil

}

func (h *LikeSubHandler) GetLikes(ctx context.Context, request *pb.GetPostRequest) (*pb.MultipleReactionsResponse, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	likeDetails := h.service.GetLikesForPost(ctx, postId)
	reactionsResponse, ok := funk.Map(likeDetails, func(dto *domain.LikeDetailsDTO) *pb.Reaction { return mapLike(dto) }).([]*pb.Reaction)
	if !ok {
		panic(fmt.Errorf("error during conversion of posts"))
	}
	return &pb.MultipleReactionsResponse{Reactions: reactionsResponse}, nil
}

func (h *LikeSubHandler) UndoLike(ctx context.Context, request *pb.GetSubresourceRequest) (*pb.EmptyRequest, error) {
	postId, err1 := primitive.ObjectIDFromHex(request.PostId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	reactionId, err2 := primitive.ObjectIDFromHex(request.SubresourceId)
	if err2 != nil {
		panic(errors.NewInvalidArgumentError("Given reaction id is invalid."))
	}
	h.service.UndoLike(ctx, postId, reactionId)
	return &pb.EmptyRequest{}, nil
}
