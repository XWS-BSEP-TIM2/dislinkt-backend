package handler

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/post/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbPost "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"net/http"
)

type PostHandler struct {
	grpcClient  *rest.ServiceClientGrpc
	errorMapper *rest.GrpcToHttpErrorCodeMapper
}

func InitPostHandler() *PostHandler {
	client := rest.InitServiceClient(config.NewConfig())
	mapper := rest.NewGrpcToHttpErrorCodeMapper()
	return &PostHandler{grpcClient: client, errorMapper: mapper}
}

func (handler *PostHandler) Get(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	res, err := postService.GetPosts(ctx, &pbPost.EmptyRequest{})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) CreatePost(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	newPost := dto.CreatePostDto{}
	if err := ctx.BindJSON(&newPost); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newPostProto := pbPost.NewPost{OwnerId: newPost.OwnerId, Content: newPost.Content, Links: newPost.Links, ImageBase64: newPost.ImageBase64}
	res, err := postService.CreatePost(ctx, &pbPost.CreatePostRequest{NewPost: &newPostProto})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *PostHandler) GetPostById(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	id := ctx.Param("id")
	res, err := postService.GetPost(ctx, &pbPost.GetPostRequest{
		PostId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetPostComments(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	id := ctx.Param("id")
	res, err := postService.GetComments(ctx, &pbPost.GetPostRequest{
		PostId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)

}

func (handler *PostHandler) CreateComment(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	comment := pbPost.NewComment{}
	if err := ctx.BindJSON(&comment); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := postService.CreateComment(ctx, &pbPost.CreateCommentRequest{NewComment: &comment, PostId: postId})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)

}

func (handler *PostHandler) GetPostComment(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	commentId := ctx.Param("comment-id")
	res, err := postService.GetComment(ctx, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: commentId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetLikes(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	res, err := postService.GetLikes(ctx, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) LikePost(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	reaction := dto.ReactionDto{}
	if err := ctx.BindJSON(&reaction); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reactionProto := pbPost.NewReaction{OwnerId: reaction.OwnerId}
	res, err := postService.GiveLike(ctx, &pbPost.CreateReactionRequest{
		PostId:      postId,
		NewReaction: &reactionProto,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *PostHandler) GetLike(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postService.GetLike(ctx, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: likeId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) RemoveLike(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postService.UndoLike(ctx, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: likeId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetDislikes(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	res, err := postService.GetDislikes(ctx, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) DislikePost(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	reaction := pbPost.NewReaction{}
	if err := ctx.BindJSON(&reaction); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := postService.GiveDislike(ctx, &pbPost.CreateReactionRequest{
		PostId:      postId,
		NewReaction: &reaction,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *PostHandler) GetDislike(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postService.GetDislike(ctx, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: likeId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) RemoveDislike(ctx *gin.Context) {
	postService := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postService.UndoDislike(ctx, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: likeId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) handleError(ctx *gin.Context, err error) {
	s, ok := status.FromError(err)
	if ok {
		httpStatus := handler.errorMapper.MapGrpcToHttpError(s.Code())
		ctx.AbortWithError(httpStatus, err)
	} else {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
}
