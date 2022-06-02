package handler

import (
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/job_offer/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbJobOffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type JobOfferHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func (handler *JobOfferHandler) Update(ctx *gin.Context) {
	jobOfferService := handler.grpcClient.JobOfferClient
	jobOfferDto := dto.JobOfferDto{}
	if err := ctx.BindJSON(&jobOfferDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jobOfferPb := pbJobOffer.JobOffer{Id: jobOfferDto.Id, UserId: jobOfferDto.UserId, CompanyName: jobOfferDto.CompanyName, Technologies: jobOfferDto.Technologies, Description: jobOfferDto.Description, Seniority: jobOfferDto.Seniority, Position: jobOfferDto.Position}
	res, _ := jobOfferService.UpdateJobOffer(ctx, &pbJobOffer.CreateJobOfferRequest{JobOffer: &jobOfferPb})
	ctx.JSON(http.StatusCreated, &res)
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

func (handler *JobOfferHandler) Create(ctx *gin.Context) {
	jobOfferService := handler.grpcClient.JobOfferClient
	jobOffer := pbJobOffer.JobOffer{}
	if err := ctx.BindJSON(&jobOffer); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, _ := jobOfferService.CreateJobOffer(ctx, &pbJobOffer.CreateJobOfferRequest{JobOffer: &jobOffer})
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *JobOfferHandler) CreateFromExternalApp(ctx *gin.Context) {
	jobOfferService := handler.grpcClient.JobOfferClient
	authService := handler.grpcClient.AuthClient

	jobOfferDto := dto.JobOfferDto{}
	if err := ctx.BindJSON(&jobOfferDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	token, err := authService.GetApiToken(ctx, &pbAuth.GetApiTokenRequest{TokenCode: jobOfferDto.ApiToken})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	jobOffer := pbJobOffer.JobOffer{
		Id:           jobOfferDto.Id,
		UserId:       token.Token.UserId,
		CompanyName:  jobOfferDto.CompanyName,
		Technologies: jobOfferDto.Technologies,
		Description:  jobOfferDto.Description,
		Seniority:    jobOfferDto.Seniority,
		Position:     jobOfferDto.Position,
	}

	res, _ := jobOfferService.CreateJobOffer(ctx, &pbJobOffer.CreateJobOfferRequest{JobOffer: &jobOffer})
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *JobOfferHandler) GetUserJobOffers(ctx *gin.Context) {
	id := ctx.Param("id")
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetUserJobOffers(ctx, &pbJobOffer.GetJobOfferRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) DeleteOffer(ctx *gin.Context) {
	id := ctx.Param("id")
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.DeleteJobOffer(ctx, &pbJobOffer.GetJobOfferRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func InitJobOfferHandler() *JobOfferHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &JobOfferHandler{grpcClient: client}
}
