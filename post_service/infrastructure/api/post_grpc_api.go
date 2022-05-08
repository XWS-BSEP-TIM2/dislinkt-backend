package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/ecoding"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/errors"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service           *application.PostService
	commentSubHandler *CommentSubHandler
	likeSubHandler    *LikeSubHandler
}

func NewPostHandler(service *application.PostService) *PostHandler {
	return &PostHandler{
		service:           service,
		commentSubHandler: NewCommentHandler(service),
		likeSubHandler:    NewLikeHandler(service),
	}
}

func (handler *PostHandler) GetPost(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.PostResponse, err error) {
	defer handleError(&err)
	objectId, idFromHexErr := primitive.ObjectIDFromHex(request.PostId)
	if idFromHexErr != nil {
		return nil, err
	}
	postDetails := handler.service.GetPost(ctx, objectId)
	postResponse = mapPostDetailsToResponse(postDetails)
	err = nil
	return
}

func (handler *PostHandler) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (postResponse *pb.PostResponse, err error) {
	defer handleError(&err)
	postDetails := handler.service.CreatePost(ctx, mapNewPostToPost(request.NewPost))
	postResponse = mapPostDetailsToResponse(postDetails)
	err = nil
	return
}

func (handler *PostHandler) GetPosts(ctx context.Context, request *pb.EmptyRequest) (postResponse *pb.MultiplePostsResponse, err error) {
	defer handleError(&err)
	postsDetails := handler.service.GetPosts(ctx)
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

func handleError(err *error) {
	if r := recover(); r != nil {
		e, ok := r.(error)
		if ok {
			*err = e
		}
	}
}

// comments subresource

func (handler *PostHandler) GetComment(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.CommentResponse, err error) {
	defer handleError(&err)
	return handler.commentSubHandler.GetComment(ctx, request)
}

func (handler *PostHandler) CreateComment(ctx context.Context, request *pb.CreateCommentRequest) (postResponse *pb.CommentResponse, err error) {
	defer handleError(&err)
	return handler.commentSubHandler.CreateComment(ctx, request)
}

func (handler *PostHandler) GetComments(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.MultipleCommentsResponse, err error) {
	defer handleError(&err)
	return handler.commentSubHandler.GetComments(ctx, request)
}

// likes subresource

func (handler *PostHandler) GetLike(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.ReactionResponse, err error) {
	defer handleError(&err)
	return handler.likeSubHandler.GetLike(ctx, request)
}

func (handler *PostHandler) GiveLike(ctx context.Context, request *pb.CreateReactionRequest) (postResponse *pb.ReactionResponse, err error) {
	defer handleError(&err)
	return handler.likeSubHandler.GiveLike(ctx, request)
}

func (handler *PostHandler) GetLikes(ctx context.Context, request *pb.GetPostRequest) (postResponse *pb.MultipleReactionsResponse, err error) {
	defer handleError(&err)
	return handler.likeSubHandler.GetLikes(ctx, request)
}

func (handler *PostHandler) UndoLike(ctx context.Context, request *pb.GetSubresourceRequest) (postResponse *pb.EmptyRequest, err error) {
	defer handleError(&err)
	return handler.likeSubHandler.UndoLike(ctx, request)
}
