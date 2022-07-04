package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/ecoding"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/api/error_mappers"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service           *application.PostService
	commentSubHandler *CommentSubHandler
	likeSubHandler    *LikeSubHandler
	dislikeSubHandler *DislikeSubHandler
	errorMapper       *error_mappers.ErrorMapperRegistry
}

func NewPostHandler(service *application.PostService) *PostHandler {
	return &PostHandler{
		service:           service,
		commentSubHandler: NewCommentHandler(service),
		likeSubHandler:    NewLikeHandler(service),
		dislikeSubHandler: NewDislikeHandler(service),
		errorMapper:       error_mappers.NewErrorMapperRegistry(),
	}
}

func (handler *PostHandler) GetPost(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.PostResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetPost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	objectId, idFromHexErr := primitive.ObjectIDFromHex(request.PostId)
	if idFromHexErr != nil {
		return nil, err
	}
	postDetails := handler.service.GetPost(ctx2, objectId)
	postResponse = mapPostDetailsToResponse(postDetails)
	err = nil
	return
}

func (handler *PostHandler) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (postResponse *pb.PostResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "CreatePost")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	postDetails := handler.service.CreatePost(ctx2, mapNewPostToPost(request.NewPost))
	postResponse = mapPostDetailsToResponse(postDetails)
	err = nil
	return
}

func (handler *PostHandler) GetPosts(ctx context.Context, request *pb.EmptyRequest) (postResponse *pb.MultiplePostsResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetPosts")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	postsDetails := handler.service.GetPosts(ctx2)
	postsResponse, ok := funk.Map(postsDetails, func(dto *domain.PostDetailsDTO) *pb.Post { return mapPost(dto) }).([]*pb.Post)
	if !ok {
		panic(fmt.Errorf("error during conversion of posts"))
	}
	postResponse = &pb.MultiplePostsResponse{Posts: postsResponse}
	err = nil
	return
}

func (handler *PostHandler) GetPostsFromUser(ctx context.Context, request *pb.GetPostsFromUserRequest) (postResponse *pb.MultiplePostsResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetPostsFromUser")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	userId, idFromHexErr := primitive.ObjectIDFromHex(request.UserId)
	if idFromHexErr != nil {
		return nil, err
	}
	postsDetails := handler.service.GetPostsFromUser(ctx2, userId)
	postsResponse, ok := funk.Map(postsDetails, func(dto *domain.PostDetailsDTO) *pb.Post { return mapPost(dto) }).([]*pb.Post)
	if !ok {
		panic(fmt.Errorf("error during conversion of posts"))
	}
	postResponse = &pb.MultiplePostsResponse{Posts: postsResponse}
	err = nil
	return
}

func mapPostDetailsToResponse(postDetails *domain.PostDetailsDTO) *pb.PostResponse {
	return &pb.PostResponse{Post: mapPost(postDetails)}
}

func mapNewPostToPost(newPost *pb.NewPost) *domain.Post {
	ownerId, err1 := primitive.ObjectIDFromHex(newPost.OwnerId)
	if err1 != nil {
		panic(errors.NewInvalidArgumentError("Given post id is invalid."))
	}
	imageBytes, err2 := ecoding.NewBase64Coder().Decode(newPost.ImageBase64)
	if err2 != nil {
		panic(fmt.Errorf("error during conversion of image"))
	}
	return &domain.Post{
		OwnerId:      ownerId,
		CreationTime: time.Now(),
		Content:      newPost.Content,
		Image:        imageBytes,
		Links:        newPost.Links,
		Comments:     []*domain.Comment{},
		Likes:        []*domain.Like{},
		Dislikes:     []*domain.Dislike{},
	}
}

func (handler *PostHandler) handleError(err *error) {
	if r := recover(); r != nil {
		e, ok := r.(error)
		if ok {
			*err = handler.errorMapper.ToStatusError(e)
		}
	}
}

// comments subresource

func (handler *PostHandler) GetComment(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.CommentResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetComment")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.commentSubHandler.GetComment(ctx2, request)
}

func (handler *PostHandler) CreateComment(ctx context.Context, request *pb.CreateCommentRequest) (postResponse *pb.CommentResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "CreateComment")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.commentSubHandler.CreateComment(ctx2, request)
}

func (handler *PostHandler) GetComments(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.MultipleCommentsResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetComments")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.commentSubHandler.GetComments(ctx2, request)
}

// likes subresource

func (handler *PostHandler) GetLike(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.ReactionResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetLike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.likeSubHandler.GetLike(ctx2, request)
}

func (handler *PostHandler) GiveLike(ctx context.Context, request *pb.CreateReactionRequest) (postResponse *pb.ReactionResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GiveLike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.likeSubHandler.GiveLike(ctx2, request)
}

func (handler *PostHandler) GetLikes(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.MultipleReactionsResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetLikes")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.likeSubHandler.GetLikes(ctx2, request)
}

func (handler *PostHandler) UndoLike(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.EmptyRequest, err error) {
	span := tracer.StartSpanFromContext(ctx, "UndoLike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.likeSubHandler.UndoLike(ctx2, request)
}

// dislikes subresource

func (handler *PostHandler) GetDislike(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.ReactionResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.dislikeSubHandler.GetDislike(ctx2, request)
}

func (handler *PostHandler) GiveDislike(ctx context.Context, request *pb.CreateReactionRequest) (postResponse *pb.ReactionResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GiveDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.dislikeSubHandler.GiveDislike(ctx2, request)
}

func (handler *PostHandler) GetDislikes(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.MultipleReactionsResponse, err error) {
	span := tracer.StartSpanFromContext(ctx, "GetDislikes")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.dislikeSubHandler.GetDislikes(ctx2, request)
}

func (handler *PostHandler) UndoDislike(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.EmptyRequest, err error) {
	span := tracer.StartSpanFromContext(ctx, "UndoDislike")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(ctx, span)

	defer handler.handleError(&err)
	return handler.dislikeSubHandler.UndoDislike(ctx2, request)
}
