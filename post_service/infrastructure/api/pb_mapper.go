package api

import (
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
)

func mapProduct(post *domain.Post) *pb.Post {
	postPb := &pb.Post{
		Id:   post.Id.Hex(),
		Name: post.Name,
	}

	return postPb
}
