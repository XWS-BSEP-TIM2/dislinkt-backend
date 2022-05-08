package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain/ecoding"
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service *application.PostService
}

func NewPostHandler(service *application.PostService) *PostHandler {
	return &PostHandler{
		service: service,
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
		panic(fmt.Errorf("error during conversion of id: %s", newPost.OwnerId))
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
		Comments:     []domain.Comment{},
		Likes:        []domain.Like{},
		Dislikes:     []domain.Dislike{},
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
