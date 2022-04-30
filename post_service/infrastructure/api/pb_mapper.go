package api

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
)

func mapProduct(product *domain.Post) *pb.Post {
	postPb := &pb.Post{
		Id:   product.Id.Hex(),
		Name: product.Name,
	}

	return postPb
}
