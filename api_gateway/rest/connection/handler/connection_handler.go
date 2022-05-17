package handler

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/DTO"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConnectionHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func InitConnectionHandler() *ConnectionHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &ConnectionHandler{grpcClient: client}
}

func (handler *ConnectionHandler) GetBlockedUsers(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	userId := ctx.Param("id")
	res, err := connectionService.GetBlockeds(ctx, &pbConnection.GetRequest{
		UserID: userId,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) AddFriend(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	addFriendDto := pbConnection.UserAction{}
	if err := ctx.BindJSON(&addFriendDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.AddFriend(ctx, &pbConnection.AddFriendRequest{
		AddFriendDTO: &addFriendDto,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) RemoveFriend(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	removeFriendDto := pbConnection.UserAction{}
	if err := ctx.BindJSON(&removeFriendDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.RemoveFriend(ctx, &pbConnection.RemoveFriendRequest{
		RemoveFriendDTO: &removeFriendDto,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) BlockUser(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	blockUser := pbConnection.UserAction{}
	if err := ctx.BindJSON(&blockUser); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.AddBlockUser(ctx, &pbConnection.AddBlockUserRequest{
		AddBlockUserDTO: &blockUser,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) UnblockUser(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	unblockUser := pbConnection.UserAction{}
	if err := ctx.BindJSON(&unblockUser); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.UnblockUser(ctx, &pbConnection.UnblockUserRequest{
		UnblockUserDTO: &unblockUser,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) CreateFriendRequest(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	friendRequest := pbConnection.UserAction{}
	if err := ctx.BindJSON(&friendRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.SendFriendRequest(ctx, &pbConnection.SendFriendRequestRequest{
		SendFriendRequestRequestDTO: &friendRequest,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) DeleteFriendRequest(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	friendRequest := pbConnection.UserAction{}
	if err := ctx.BindJSON(&friendRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := connectionService.UnsendFriendRequest(ctx, &pbConnection.UnsendFriendRequestRequest{
		UnsendFriendRequestRequestDTO: &friendRequest,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *ConnectionHandler) GetConnectionDetails(ctx *gin.Context) {
	connectionService := handler.grpcClient.ConnectionClient
	userIda := ctx.Param("ida")
	userIdb := ctx.Param("idb")
	res, err := connectionService.GetConnectionDetail(ctx, &pbConnection.GetConnectionDetailRequest{
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

	connectionService := handler.grpcClient.ConnectionClient
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	recommendations, err := connectionService.GetBlockeds(ctx, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	recommendationDTO := handler.GenerateProfileDTO(ctx, recommendations.Users)
	ctx.JSON(http.StatusOK, &recommendationDTO)
}

func (handler *ConnectionHandler) GetFriendRequests(ctx *gin.Context) {
	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	connectionService := handler.grpcClient.ConnectionClient
	recommendations, err := connectionService.GetFriendRequests(ctx, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx, recommendations.Users)
	ctx.JSON(http.StatusOK, &recommendationDTO)
}

func (handler *ConnectionHandler) GetRecommendation(ctx *gin.Context) {

	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	connectionService := handler.grpcClient.ConnectionClient
	recommendations, err := connectionService.GetRecommendation(ctx, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx, recommendations.Users)

	ctx.JSON(http.StatusOK, &recommendationDTO)
}

func (handler *ConnectionHandler) GetFriends(ctx *gin.Context) {

	userId := ctx.Param("id")

	connectionService := handler.grpcClient.ConnectionClient
	friends, err := connectionService.GetFriends(ctx, &pbConnection.GetRequest{UserID: userId})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	friendsDTO := handler.GenerateProfileDTO(ctx, friends.Users)

	ctx.JSON(http.StatusOK, &friendsDTO)
}

func (handler *ConnectionHandler) GenerateProfileDTO(ctx context.Context, connUsers []*pbConnection.User) []*DTO.ConnectionDTO {
	profileService := handler.grpcClient.ProfileClient

	var usersDTO []*DTO.ConnectionDTO

	for _, v := range connUsers {
		profile, err := profileService.Get(ctx, &pbProfile.GetRequest{Id: v.UserID})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		usersDTO = append(usersDTO, DTO.MapConnectionDTO(profile.Profile))
	}
	return usersDTO
}
