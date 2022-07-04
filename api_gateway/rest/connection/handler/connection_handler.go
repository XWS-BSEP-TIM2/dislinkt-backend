package handler

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type ConnectionHandler struct {
	grpcClient *rest.ServiceClientGrpc
	tracer     opentracing.Tracer
}

func InitConnectionHandler(tracer opentracing.Tracer) *ConnectionHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &ConnectionHandler{grpcClient: client, tracer: tracer}
}

func (handler *ConnectionHandler) GetBlockedUsers(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetBlockedUsers", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	userId := ctx.Param("id")
	res, err := connectionService.GetBlockeds(ctx2, &pbConnection.GetRequest{
		UserID: userId,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) AddFriend(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("AddFriend", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	addFriendDto := pbConnection.UserAction{}
	if err := ctx.BindJSON(&addFriendDto); err != nil {
		handler.loggError(ctx2, "AddFriend", "", ctx.ClientIP(), err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.AddFriend(ctx2, &pbConnection.AddFriendRequest{
		AddFriendDTO: &addFriendDto,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) RemoveFriend(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("RemoveFriend", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	removeFriendDto := pbConnection.UserAction{}
	if err := ctx.BindJSON(&removeFriendDto); err != nil {
		handler.loggError(ctx2, "RemoveFriend", "", ctx.ClientIP(), err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.RemoveFriend(ctx2, &pbConnection.RemoveFriendRequest{
		RemoveFriendDTO: &removeFriendDto,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) BlockUser(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("BlockUser", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	blockUser := pbConnection.UserAction{}
	if err := ctx.BindJSON(&blockUser); err != nil {
		handler.loggError(ctx2, "BlockUser", "", ctx.ClientIP(), err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.AddBlockUser(ctx2, &pbConnection.AddBlockUserRequest{
		AddBlockUserDTO: &blockUser,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) UnblockUser(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("UnblockUser", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	unblockUser := pbConnection.UserAction{}
	if err := ctx.BindJSON(&unblockUser); err != nil {
		handler.loggError(ctx2, "UnblockUser", "", ctx.ClientIP(), err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.UnblockUser(ctx2, &pbConnection.UnblockUserRequest{
		UnblockUserDTO: &unblockUser,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) CreateFriendRequest(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("CreateFriendRequest", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	friendRequest := pbConnection.UserAction{}
	if err := ctx.BindJSON(&friendRequest); err != nil {
		handler.loggError(ctx2, "CreateFriendRequest", "", ctx.ClientIP(), err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.SendFriendRequest(ctx2, &pbConnection.SendFriendRequestRequest{
		SendFriendRequestRequestDTO: &friendRequest,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) DeleteFriendRequest(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("DeleteFriendRequest", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	friendRequest := pbConnection.UserAction{}
	if err := ctx.BindJSON(&friendRequest); err != nil {
		handler.loggError(ctx2, "DeleteFriendRequest", "", ctx.ClientIP(), err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.UnsendFriendRequest(ctx2, &pbConnection.UnsendFriendRequestRequest{
		UnsendFriendRequestRequestDTO: &friendRequest,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) GetConnectionDetails(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetConnectionDetails", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	userIda := ctx.Param("ida")
	userIdb := ctx.Param("idb")
	res, err := connectionService.GetConnectionDetail(ctx2, &pbConnection.GetConnectionDetailRequest{
		UserIDa: userIda,
		UserIDb: userIdb,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) GetBlocks(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetBlocks", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	connectionService := handler.grpcClient.ConnectionClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	recommendations, err := connectionService.GetBlockeds(ctx2, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	recommendationDTO := handler.GenerateProfileDTO(ctx2, recommendations.Users)
	ctx.JSON(http.StatusOK, &recommendationDTO)
}

func (handler *ConnectionHandler) GetFriendRequests(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetFriendRequests", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	connectionService := handler.grpcClient.ConnectionClient
	recommendations, err := connectionService.GetFriendRequests(ctx2, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx2, recommendations.Users)
	ctx.JSON(http.StatusOK, &recommendationDTO)
}

func (handler *ConnectionHandler) GetRecommendation(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetRecommendation", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	connectionService := handler.grpcClient.ConnectionClient
	recommendations, err := connectionService.GetRecommendation(ctx2, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx2, recommendations.Users)

	ctx.JSON(http.StatusOK, &recommendationDTO)
}

func (handler *ConnectionHandler) GetFriends(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetFriends", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := ctx.Param("id")

	connectionService := handler.grpcClient.ConnectionClient
	friends, err := connectionService.GetFriends(ctx2, &pbConnection.GetRequest{UserID: userId})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	friendsDTO := handler.GenerateProfileDTO(ctx2, friends.Users)

	ctx.JSON(http.StatusOK, &friendsDTO)
}

func (handler *ConnectionHandler) GenerateProfileDTO(ctx context.Context, connUsers []*pbConnection.User) []*ConnectionDTO {
	span := tracer.StartSpanFromContext(ctx, "GenerateProfileDTO")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	profileService := handler.grpcClient.ProfileClient

	var usersDTO []*ConnectionDTO

	for _, v := range connUsers {
		profile, err := profileService.Get(ctx2, &pbProfile.GetRequest{Id: v.UserID})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		usersDTO = append(usersDTO, MapConnectionDTO(profile.Profile))
	}
	return usersDTO
}

func (handler *ConnectionHandler) loggError(ctx context.Context, serviceFunctionName, userID, ipAddress, description string) {
	span := tracer.StartSpanFromContext(ctx, "loggError")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	handler.grpcClient.LoggingClient.LoggError(ctx2, &loggingS.LogRequest{
		ServiceName:         "API GATEWAY",
		ServiceFunctionName: serviceFunctionName,
		UserID:              userID,
		IpAddress:           ipAddress,
		Description:         description,
	})
}
