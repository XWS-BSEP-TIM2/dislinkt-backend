package api

import (
	"context"
	"github.com/tamararankovic/microservices_demo/catalogue_service/application"
	pb "github.com/tamararankovic/microservices_demo/common/proto/post_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostHandler struct {
	pb.UnimplementedPostServiceServer
	service *application.PostService
}

func NewProductHandler(service *application.PostService) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (handler *PostHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	post, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	postPb := mapProduct(post)
	response := &pb.GetResponse{
		Post: postPb,
	}
	return response, nil
}
