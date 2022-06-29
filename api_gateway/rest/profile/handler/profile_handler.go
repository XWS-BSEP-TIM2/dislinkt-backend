package handler

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbJobOffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	tracer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/validators"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"io"
	"net/http"
)

type ProfileHandler struct {
	grpcClient *rest.ServiceClientGrpc
	tracer     opentracing.Tracer
	closer     io.Closer
}

func InitProfileHandler() *ProfileHandler {
	tracer, closer := tracer.Init("profile_service")
	opentracing.SetGlobalTracer(tracer)
	client := rest.InitServiceClient(config.NewConfig())
	return &ProfileHandler{grpcClient: client,
		tracer: tracer,
		closer: closer}
}

func (handler *ProfileHandler) Get(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Get", handler.tracer, ctx.Request)
	defer span.Finish()

	ctx2 := tracer.ContextWithSpan(context.Background(), span)
	profileService := handler.grpcClient.ProfileClient
	res, err := profileService.GetAll(ctx2, &pbProfile.EmptyRequest{})
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
	connectionService := handler.grpcClient.ConnectionClient
	authService := handler.grpcClient.AuthClient

	profile := pbProfile.Profile{}

	if err := ctx.BindJSON(&profile); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	authService.EditData(ctx, &pbAuth.EditDataRequest{Email: profile.Email, Username: profile.Username, UserId: profile.Id, IsTwoFactor: profile.IsTwoFactor})
	res, err1 := profileService.UpdateProfile(ctx, &pbProfile.CreateProfileRequest{Profile: &profile})
	if err1 != nil {
		ctx.AbortWithError(http.StatusBadGateway, err1)
	}
	cpb := &pbConnection.ChangePrivacyBody{UserID: profile.Id, IsPrivate: profile.IsPrivate}
	_, err2 := connectionService.ChangePrivacy(ctx, &pbConnection.ChangePrivacyRequest{ChangePrivacyBody: cpb})
	if err2 != nil {
		ctx.AbortWithError(http.StatusBadGateway, err2)
	}
	errUpdateSkillsInJobOfferService := handler.updateSkillsInJobOfferService(ctx, &profile)
	if errUpdateSkillsInJobOfferService != nil {
		ctx.AbortWithError(http.StatusBadGateway, errUpdateSkillsInJobOfferService)
	}

	ctx.JSON(http.StatusCreated, &res)
}

func (handler *ProfileHandler) updateSkillsInJobOfferService(ctx *gin.Context, profile *pbProfile.Profile) error {
	jobOfferService := handler.grpcClient.JobOfferClient
	var technologies []string
	addTechnologie := true
	for _, s1 := range profile.Skills {
		addTechnologie = true
		for _, t := range technologies {
			if s1.Name == t {
				addTechnologie = false
				break
			}
		}
		if addTechnologie {
			technologies = append(technologies, s1.Name)
		}
	}

	registrationResult, err := jobOfferService.UpdateUserSkills(ctx, &pbJobOffer.UpdateUserSkillsRequest{UserID: profile.Id, Technologies: technologies})
	fmt.Println(registrationResult)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (handler *ProfileHandler) ChangePassword(ctx *gin.Context) {
	authService := handler.grpcClient.AuthClient

	changePasswordRequest := pbAuth.ChangePasswordRequest{}

	if err := ctx.BindJSON(&changePasswordRequest); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if !validators.IsPasswordCracked(changePasswordRequest.NewPassword) {
		ctx.JSON(http.StatusUnprocessableEntity, pbAuth.ChangePasswordResponse{Status: http.StatusUnprocessableEntity, Msg: "Password is already cracked"})
		return
	}

	res, err := authService.ChangePassword(ctx, &changePasswordRequest)
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
