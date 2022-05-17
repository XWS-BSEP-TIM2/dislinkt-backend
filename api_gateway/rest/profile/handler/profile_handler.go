package handler

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfileHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func InitProfileHandler() *ProfileHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &ProfileHandler{grpcClient: client}
}

func (handler *ProfileHandler) Get(ctx *gin.Context) {
	profileService := handler.grpcClient.ProfileClient
	res, err := profileService.GetAll(ctx, &pbProfile.EmptyRequest{})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusCreated, &res)

}

func (handler *ProfileHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	profileService := handler.grpcClient.ProfileClient
	res, err := profileService.Get(ctx, &pbProfile.GetRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *ProfileHandler) Update(ctx *gin.Context) {
	profileService := handler.grpcClient.ProfileClient

	profile := pbProfile.Profile{}

	if err := ctx.BindJSON(&profile); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := profileService.UpdateProfile(ctx, &pbProfile.CreateProfileRequest{Profile: &profile})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *ProfileHandler) Search(ctx *gin.Context) {
	profileService := handler.grpcClient.ProfileClient

	searchRequest := pbProfile.SearchProfileRequest{}

	if err := ctx.BindJSON(&searchRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := profileService.SearchProfile(ctx, &searchRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusCreated, &res)
}
