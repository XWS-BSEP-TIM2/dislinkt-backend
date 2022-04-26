package api

import (
	"github.com/tamararankovic/microservices_demo/catalogue_service/domain"
	pb "github.com/tamararankovic/microservices_demo/common/proto/post_service"
)

func mapProduct(product *domain.Post) *pb.Post {
	postPb := &pb.Post{
		Id:   product.Id.Hex(),
		Name: product.Name,
	}

	return postPb
}
