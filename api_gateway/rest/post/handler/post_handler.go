package handler

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/post/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbPost "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

type PostHandler struct {
	grpcClient  *rest.ServiceClientGrpc
	errorMapper *GrpcToHttpErrorCodeMapper
}

func InitPostHandler() *PostHandler {
	client := rest.InitServiceClient(config.NewConfig())
	mapper := NewGrpcToHttpErrorCodeMapper()
	return &PostHandler{grpcClient: client, errorMapper: mapper}
}

func (handler *PostHandler) Get(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	res, err := postClient.GetPosts(ctxt, &pbPost.EmptyRequest{})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) CreatePost(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	newPost := dto.CreatePostDto{}
	if err := ctx.BindJSON(&newPost); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newPostProto := pbPost.NewPost{OwnerId: newPost.OwnerId, Content: newPost.Content, Links: newPost.Links, ImageBase64: newPost.ImageBase64}
	res, err := postClient.CreatePost(ctxt, &pbPost.CreatePostRequest{NewPost: &newPostProto})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *PostHandler) GetPostById(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	id := ctx.Param("id")
	res, err := postClient.GetPost(ctxt, &pbPost.GetPostRequest{
		PostId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetPostsFromUser(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	id := ctx.Param("user-id")
	res, err := postClient.GetPostsFromUser(ctxt, &pbPost.GetPostsFromUserRequest{
		UserId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetPostComments(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	id := ctx.Param("id")
	res, err := postClient.GetComments(ctxt, &pbPost.GetPostRequest{
		PostId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)

}

func (handler *PostHandler) CreateComment(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	comment := dto.CreateCommentDto{}
	if err := ctx.BindJSON(&comment); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := postClient.CreateComment(ctxt, &pbPost.CreateCommentRequest{
		NewComment: &pbPost.NewComment{
			OwnerId: comment.OwnerId,
			Content: comment.Content,
		},
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)

}

func (handler *PostHandler) GetPostComment(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	commentId := ctx.Param("comment-id")
	res, err := postClient.GetComment(ctxt, &pbPost.GetSubresourceRequest{
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
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	res, err := postClient.GetLikes(ctxt, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) LikePost(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	reaction := dto.ReactionDto{}
	if err := ctx.BindJSON(&reaction); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reactionProto := pbPost.NewReaction{OwnerId: reaction.OwnerId}
	res, err := postClient.GiveLike(ctxt, &pbPost.CreateReactionRequest{
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
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postClient.GetLike(ctxt, &pbPost.GetSubresourceRequest{
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
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postClient.UndoLike(ctxt, &pbPost.GetSubresourceRequest{
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
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	res, err := postClient.GetDislikes(ctxt, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) DislikePost(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	reaction := dto.ReactionDto{}
	if err := ctx.BindJSON(&reaction); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	reactionProto := pbPost.NewReaction{OwnerId: reaction.OwnerId}
	res, err := postClient.GiveDislike(ctxt, &pbPost.CreateReactionRequest{
		PostId:      postId,
		NewReaction: &reactionProto,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *PostHandler) GetDislike(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	dislikeId := ctx.Param("dislike-id")
	res, err := postClient.GetDislike(ctxt, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: dislikeId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) RemoveDislike(ctx *gin.Context) {
	ctxt := handler.appendTokenToContext(ctx)
	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	dislikeId := ctx.Param("dislike-id")
	res, err := postClient.UndoDislike(ctxt, &pbPost.GetSubresourceRequest{
		PostId:        postId,
		SubresourceId: dislikeId,
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

func (handler *PostHandler) appendTokenToContext(ctx *gin.Context) context.Context {
	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		return ctx
	} else {
		return metadata.AppendToOutgoingContext(ctx, "authorization", auth)
	}
}
