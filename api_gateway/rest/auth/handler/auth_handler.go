package handler

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/auth/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbJobOffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthHandler struct {
	grpcClient *rest.ServiceClientGrpc
	tracer     opentracing.Tracer
}

func InitAuthHandler(tracer opentracing.Tracer) *AuthHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &AuthHandler{
		grpcClient: client,
		tracer:     tracer,
	}
}
func (authHandler *AuthHandler) Login(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Login", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	loginDto := dto.LoginRequestDto{}

	if err := ctx.BindJSON(&loginDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	v := validator.New()
	validators.UsernameValidator(ctx2, v)
	errV := v.Struct(loginDto)
	if errV != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
			Message: "Description is not valid",
		})
		return
	}

	res, err := authHandler.grpcClient.AuthClient.Login(ctx2, &pbAuth.LoginRequest{
		Username: loginDto.Username,
		Password: loginDto.Password,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) Verify(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Verify", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	username := ctx.Param("username")
	code := ctx.Param("code")
	v := pbAuth.VerifyRequest{Username: username, Code: code}

	res, err := authHandler.grpcClient.AuthClient.Verify(ctx2, &v)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) ResendVerify(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("ResendVerify", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	username := ctx.Param("username")
	v := pbAuth.ResendVerifyRequest{Username: username}

	res, err := authHandler.grpcClient.AuthClient.ResendVerify(ctx2, &v)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) GetRecovery(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GetRecovery", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	username := ctx.Param("username")
	r := pbAuth.RecoveryRequest{Username: username}
	res, err := authHandler.grpcClient.AuthClient.Recovery(ctx2, &r)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) Recover(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Recover", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	rl := pbAuth.RecoveryRequestLogin{}

	if err := ctx.BindJSON(&rl); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fmt.Println(rl)

	if !validators.IsPasswordCracked(rl.NewPassword) {
		ctx.JSON(http.StatusUnprocessableEntity, pbAuth.LoginResponse{Status: http.StatusUnprocessableEntity, Error: "Password is already cracked"})
		return
	}

	res, err := authHandler.grpcClient.AuthClient.Recover(ctx2, &rl)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) Register(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Register", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	registerDto := dto.RegisterDTO{}

	if err := ctx.BindJSON(&registerDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	v := validator.New()
	validators.UsernameValidator(ctx2, v)
	validators.EmailValidator(ctx2, v)
	validators.PasswordValidator(ctx2, v)
	validators.NameValidator(ctx2, v)
	validators.NumberValidator(ctx2, v)
	errV := v.Struct(registerDto)
	if errV != nil {
		if strings.Contains(errV.Error(), "Position") {
			ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
				Message: "Position is not valid",
			})
		} else if strings.Contains(errV.Error(), "Seniority") {
			ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
				Message: "Seniority is not valid",
			})
		} else if strings.Contains(errV.Error(), "Description") {
			ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
				Message: "Description is not valid",
			})
		} else if strings.Contains(errV.Error(), "Email") {
			ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
				Message: "Email is not valid",
			})
		} else if strings.Contains(errV.Error(), "PhoneNumber") {
			ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
				Message: "Phone number is not valid",
			})
		} else {
			ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
				Message: "Password is not valid",
			})
		}
		return
	}

	if !validators.IsPasswordCracked(registerDto.Password) {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
			Message: "Password is already cracked",
		})
		return
	}

	userID, errAuth := authHandler.registerAuth(registerDto)
	if errAuth != nil {
		ctx.AbortWithError(http.StatusBadGateway, errAuth)
		return
	}

	//errProfile := authHandler.registerProfile(userID, registerDto)
	//if errProfile != nil {
	//	ctx.AbortWithError(http.StatusBadGateway, errProfile)
	//	return
	//}
	//
	//errConnection := authHandler.registerConnection(userID, registerDto.IsPrivate)
	//if errProfile != nil {
	//	ctx.AbortWithError(http.StatusBadGateway, errConnection)
	//	return
	//}

	/*
			//TODO: poziv reg usera u jobOffer graf bazi
		errRegUserInJobOffer := authHandler.registerUserInJobOffer(userID)
		if errRegUserInJobOffer != nil {
			ctx.AbortWithError(http.StatusBadGateway, errConnection)
			return
		}
	*/

	fmt.Println("successfully registered new user with ID:", userID)

	responsDTO := dto.RegisterResponsDTO{Id: userID, Username: registerDto.Username}

	ctx.JSON(http.StatusCreated, &responsDTO)
}

func (authHandler *AuthHandler) SendMailForMagicLinkRegistration(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("SendMailForMagicLinkRegistration", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	authService := authHandler.grpcClient.AuthClient
	b := pbAuth.EmailForPasswordlessLoginRequest{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	response, _ := authService.SendEmailForPasswordlessLogin(ctx2, &b)
	ctx.JSON(http.StatusOK, &response)
}

func (authHandler *AuthHandler) registerAuth(registerDTO dto.RegisterDTO) (string, error) {
	t, err := time.Parse("2022-02-25", registerDTO.Birthday)
	if err != nil {
		fmt.Println("Error BirthDate")
	}
	authS := authHandler.grpcClient.AuthClient
	response, err := authS.Register(context.TODO(), &pbAuth.RegisterRequest{Username: registerDTO.Username, Password: registerDTO.Password, Email: registerDTO.Email, Gender: registerDTO.Gender, Name: registerDTO.Name, BirthDate: timestamppb.New(t), Surname: registerDTO.Surname, PhoneNumber: registerDTO.PhoneNumber, IsPrivate: registerDTO.IsPrivate})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(response)
	return response.UserID, err
}

func (authHandler *AuthHandler) registerProfile(userID string, registerDTO dto.RegisterDTO) error {
	profileService := authHandler.grpcClient.ProfileClient
	response, err := profileService.CreateProfile(context.TODO(), &pbProfile.CreateProfileRequest{Profile: registerDTO.ToProto(userID)})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(response)
	return nil
}

func (authHandler *AuthHandler) registerConnection(userID string, IsPrivate bool) error {
	connectionService := authHandler.grpcClient.ConnectionClient
	registrationResult, err := connectionService.Register(context.TODO(), &pbConnection.RegisterRequest{User: &pbConnection.User{UserID: userID, IsPrivate: IsPrivate}})
	fmt.Println(registrationResult)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (authHandler *AuthHandler) registerUserInJobOffer(userID string) error {
	jobOfferService := authHandler.grpcClient.JobOfferClient
	registrationResult, err := jobOfferService.CreateUser(context.TODO(), &pbJobOffer.CreateUserRequest{UserID: userID})
	fmt.Println(registrationResult)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (authHandler *AuthHandler) MagicLinkLogin(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("MagicLinkLogin", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	authService := authHandler.grpcClient.AuthClient
	passwordlessMessage := pbAuth.PasswordlessLoginRequest{}
	if err := ctx.BindJSON(&passwordlessMessage); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	login, err := authService.PasswordlessLogin(ctx2, &passwordlessMessage)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &login)
		return
	}
	ctx.JSON(http.StatusOK, &login)
}

func (authHandler *AuthHandler) GenerateApiToken(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GenerateApiToken", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := ctx.Param("userId")
	authService := authHandler.grpcClient.AuthClient
	res, err := authService.GenerateApiToken(ctx2, &pbAuth.ApiTokenRequest{UserId: userId})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, &res)
		return
	}
	if res.Error != nil {
		ctx.JSON(int(res.Error.ErrorCode), &res.Error)
		return
	}
	ctx.JSON(http.StatusCreated, &res)

}

func (authHandler *AuthHandler) Test(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Test", authHandler.tracer, ctx.Request)
	defer span.Finish()

	ctx.JSON(http.StatusOK, "Everything is OK")
}

func (authHandler *AuthHandler) GenerateQrCode(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("GenerateQrCode", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userId := ctx.Param("id")
	authService := authHandler.grpcClient.AuthClient
	res, err := authService.GenerateQr2TF(ctx2, &pbAuth.UserIdRequest{UserId: userId})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, &res)
		return
	}
	if res.Error != nil {
		ctx.JSON(int(res.Error.ErrorCode), &res.Error)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}

func (authHandler *AuthHandler) Verify2Factor(ctx *gin.Context) {
	span := tracer.StartSpanFromRequest("Verify2Factor", authHandler.tracer, ctx.Request)
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	dto := dto.Verify2FactorDto{}
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	authService := authHandler.grpcClient.AuthClient
	res, err := authService.Verify2FactorCode(ctx2, &pbAuth.TFARequest{Code: strconv.Itoa(dto.Code), UserId: dto.UserId})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, &res)
		return
	}
	ctx.JSON(int(res.Status), &res)
}
