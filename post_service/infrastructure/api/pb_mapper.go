package api

import (
	"fmt"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapPost(postDTO *domain.PostDetailsDTO) *pb.Post {
	postId := postDTO.Post.Id.Hex()
	postPb := &pb.Post{
		//Owner: &pb.Owner{
		//	Username: postDTO.Owner.Username,
		//	Name:     postDTO.Owner.Name,
		//	Surname:  postDTO.Owner.Surname,
		//},
		CreationTime: timestamppb.New(postDTO.Post.CreationTime),
		Content:      postDTO.Post.Content,
		ImageBase64:  postDTO.ImageBase64,
		Links:        postDTO.Post.Links,
		Hrefs: []*pb.Href{
			{
				Rel: "self",
				Url: fmt.Sprintf("posts/%s", postId),
			},
			{
				Rel: "comments",
				Url: fmt.Sprintf("posts/%s/comments", postId),
			},
			{
				Rel: "likes",
				Url: fmt.Sprintf("posts/%s/likes", postId),
			},
			{
				Rel: "dislikes",
				Url: fmt.Sprintf("posts/%s/dislikes", postId),
			},
			{
				Rel: "owner",
				Url: fmt.Sprintf("profile/%s", postDTO.Post.OwnerId.Hex()),
			},
		},
		Stats: &pb.PostStats{
			CommentsNumber: int64(postDTO.Stats.CommentsNumber),
			LikesNumber:    int64(postDTO.Stats.LikesNumber),
			DislikesNumber: int64(postDTO.Stats.DislikesNumber),
		},
	}

	return postPb
}

func mapComment(commentDTO *domain.CommentDetailsDTO) *pb.Comment {
	commentPb := &pb.Comment{
		CreationTime: timestamppb.New(commentDTO.Comment.CreationTime),
		Content:      commentDTO.Comment.Content,
		Hrefs: []*pb.Href{
			{
				Rel: "self",
				Url: fmt.Sprintf("posts/%s/comments/%s", commentDTO.PostId.Hex(), commentDTO.Comment.Id.Hex()),
			},
			{
				Rel: "owner",
				Url: fmt.Sprintf("profile/%s", commentDTO.Comment.OwnerId.Hex()),
			},
		},
	}
	return commentPb
}
