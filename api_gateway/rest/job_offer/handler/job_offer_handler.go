package handler

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbJobOffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JobOfferHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func (handler *JobOfferHandler) Update(context *gin.Context) {

}

func (handler *JobOfferHandler) Get(ctx *gin.Context) {
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetAllJobOffers(ctx, &pbJobOffer.EmptyJobOfferRequest{})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetJobOffer(ctx, &pbJobOffer.GetJobOfferRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) Search(ctx *gin.Context) {
	jobOfferService := handler.grpcClient.JobOfferClient

	searchRequest := pbJobOffer.SearchJobOfferRequest{}

	if err := ctx.BindJSON(&searchRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := jobOfferService.SearchJobOffer(ctx, &searchRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusCreated, &res)
}

func InitJobOfferHandler() *JobOfferHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &JobOfferHandler{grpcClient: client}
}
