package handler

import (
	"context"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/job_offer/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/security"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbJobOffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type JobOfferHandler struct {
	grpcClient *rest.ServiceClientGrpc
	tracer     opentracing.Tracer
}

func (handler *JobOfferHandler) Update(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Update", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOfferService := handler.grpcClient.JobOfferClient
	jobOfferDto := dto.JobOfferDto{}
	if err := ctx.BindJSON(&jobOfferDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	jobOfferPb := pbJobOffer.JobOffer{Id: jobOfferDto.Id, UserId: jobOfferDto.UserId, CompanyName: jobOfferDto.CompanyName, Technologies: jobOfferDto.Technologies, Description: jobOfferDto.Description, Seniority: jobOfferDto.Seniority, Position: jobOfferDto.Position}
	res, _ := jobOfferService.UpdateJobOffer(ctx2, &pbJobOffer.CreateJobOfferRequest{JobOffer: &jobOfferPb})
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *JobOfferHandler) Get(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Get", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetAllJobOffers(ctx2, &pbJobOffer.EmptyJobOfferRequest{})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) GetById(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetById", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	id := ctx.Param("id")
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetJobOffer(ctx2, &pbJobOffer.GetJobOfferRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) GetRecommend(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetRecommend", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	dataFromToken, _ := security.ExtractTokenDataFromContext(ctx)
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetRecommendationJobOffer(ctx2, &pbJobOffer.GetRecommendationJobOfferRequest{UserID: dataFromToken.Id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) Search(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Search", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOfferService := handler.grpcClient.JobOfferClient

	searchRequest := pbJobOffer.SearchJobOfferRequest{}

	if err := ctx.BindJSON(&searchRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, err := jobOfferService.SearchJobOffer(ctx2, &searchRequest)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *JobOfferHandler) Create(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Create", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	jobOfferService := handler.grpcClient.JobOfferClient
	jobOffer := pbJobOffer.JobOffer{}
	if err := ctx.BindJSON(&jobOffer); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	res, _ := jobOfferService.CreateJobOffer(ctx2, &pbJobOffer.CreateJobOfferRequest{JobOffer: &jobOffer})
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *JobOfferHandler) CreateFromExternalApp(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("CreateFromExternalApp", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

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

	res, _ := jobOfferService.CreateJobOffer(ctx2, &pbJobOffer.CreateJobOfferRequest{JobOffer: &jobOffer})
	ctx.JSON(http.StatusCreated, &res)
}

func (handler *JobOfferHandler) GetUserJobOffers(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetUserJobOffers", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	id := ctx.Param("id")
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.GetUserJobOffers(ctx2, &pbJobOffer.GetJobOfferRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func (handler *JobOfferHandler) DeleteOffer(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("DeleteOffer", handler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	id := ctx.Param("id")
	jobOfferService := handler.grpcClient.JobOfferClient
	res, err := jobOfferService.DeleteJobOffer(ctx2, &pbJobOffer.GetJobOfferRequest{Id: id})
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
	}
	ctx.JSON(http.StatusOK, &res)
}

func InitJobOfferHandler(tracer opentracing.Tracer) *JobOfferHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &JobOfferHandler{grpcClient: client, tracer: tracer}
}
