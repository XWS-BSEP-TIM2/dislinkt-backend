package handler

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/post/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbPost "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
)

type PostHandler struct {
	grpcClient  *rest.ServiceClientGrpc
	errorMapper *GrpcToHttpErrorCodeMapper
	tracer      opentracing.Tracer
}

func InitPostHandler(tracer opentracing.Tracer) *PostHandler {
	client := rest.InitServiceClient(config.NewConfig())
	mapper := NewGrpcToHttpErrorCodeMapper()
	return &PostHandler{
		grpcClient:  client,
		errorMapper: mapper,
		tracer:      tracer,
	}
}

func (handler *PostHandler) Get(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Get", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	res, err := postClient.GetPosts(ctx2, &pbPost.EmptyRequest{})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) CreatePost(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("CreatePost", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	newPost := dto.CreatePostDto{}
	if err := ctx.BindJSON(&newPost); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newPostProto := pbPost.NewPost{OwnerId: newPost.OwnerId, Content: newPost.Content, Links: newPost.Links, ImageBase64: newPost.ImageBase64}
	res, err := postClient.CreatePost(ctx2, &pbPost.CreatePostRequest{NewPost: &newPostProto})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *PostHandler) GetPostById(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetPostById", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	id := ctx.Param("id")
	res, err := postClient.GetPost(ctx2, &pbPost.GetPostRequest{
		PostId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetPostsFromUser(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetPostsFromUser", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	id := ctx.Param("user-id")
	res, err := postClient.GetPostsFromUser(ctx2, &pbPost.GetPostsFromUserRequest{
		UserId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) GetPostComments(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetPostComments", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	id := ctx.Param("id")
	res, err := postClient.GetComments(ctx2, &pbPost.GetPostRequest{
		PostId: id,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)

}

func (handler *PostHandler) CreateComment(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("CreateComment", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	comment := dto.CreateCommentDto{}
	if err := ctx.BindJSON(&comment); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := postClient.CreateComment(ctx2, &pbPost.CreateCommentRequest{
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
	span := tracer.StartSpanFromRequest("GetPostComment", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	commentId := ctx.Param("comment-id")
	res, err := postClient.GetComment(ctx2, &pbPost.GetSubresourceRequest{
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
	span := tracer.StartSpanFromRequest("GetLikes", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	res, err := postClient.GetLikes(ctx2, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) LikePost(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("LikePost", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	reaction := dto.ReactionDto{}
	if err := ctx.BindJSON(&reaction); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reactionProto := pbPost.NewReaction{OwnerId: reaction.OwnerId}
	res, err := postClient.GiveLike(ctx2, &pbPost.CreateReactionRequest{
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
	span := tracer.StartSpanFromRequest("GetLike", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postClient.GetLike(ctx2, &pbPost.GetSubresourceRequest{
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
	span := tracer.StartSpanFromRequest("RemoveLike", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	likeId := ctx.Param("like-id")
	res, err := postClient.UndoLike(ctx2, &pbPost.GetSubresourceRequest{
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
	span := tracer.StartSpanFromRequest("GetDislikes", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	res, err := postClient.GetDislikes(ctx2, &pbPost.GetPostRequest{
		PostId: postId,
	})
	if err != nil {
		handler.handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *PostHandler) DislikePost(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("DislikePost", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	reaction := dto.ReactionDto{}
	if err := ctx.BindJSON(&reaction); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	reactionProto := pbPost.NewReaction{OwnerId: reaction.OwnerId}
	res, err := postClient.GiveDislike(ctx2, &pbPost.CreateReactionRequest{
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
	span := tracer.StartSpanFromRequest("GetDislike", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	dislikeId := ctx.Param("dislike-id")
	res, err := postClient.GetDislike(ctx2, &pbPost.GetSubresourceRequest{
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
	span := tracer.StartSpanFromRequest("RemoveDislike", handler.tracer, ctx.Request)
	defer span.Finish()
	ctxt := handler.appendTokenToContext(ctx)
	ctx2 := tracer.ContextWithSpan(ctxt, span)

	postClient := handler.grpcClient.PostClient
	postId := ctx.Param("id")
	dislikeId := ctx.Param("dislike-id")
	res, err := postClient.UndoDislike(ctx2, &pbPost.GetSubresourceRequest{
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
