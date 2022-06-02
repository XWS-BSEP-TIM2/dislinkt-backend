package handler

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/rest/auth/dto"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type AuthHandler struct {
	grpcClient *rest.ServiceClientGrpc
}

func InitAuthHandler() *AuthHandler {
	client := rest.InitServiceClient(config.NewConfig())
	return &AuthHandler{grpcClient: client}
}
func (authHandler *AuthHandler) Login(ctx *gin.Context) {

	loginDto := dto.LoginRequestDto{}

	if err := ctx.BindJSON(&loginDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	v := validator.New()
	validators.UsernameValidator(ctx, v)
	errV := v.Struct(loginDto)
	if errV != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Error{
			Message: "Description is not valid",
		})
		return
	}

	res, err := authHandler.grpcClient.AuthClient.Login(context.Background(), &pbAuth.LoginRequest{
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

	username := ctx.Param("username")
	code := ctx.Param("code")
	v := pbAuth.VerifyRequest{Username: username, Code: code}

	res, err := authHandler.grpcClient.AuthClient.Verify(context.Background(), &v)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) ResendVerify(ctx *gin.Context) {

	username := ctx.Param("username")
	v := pbAuth.ResendVerifyRequest{Username: username}

	res, err := authHandler.grpcClient.AuthClient.ResendVerify(context.Background(), &v)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) GetRecovery(ctx *gin.Context) {

	username := ctx.Param("username")
	r := pbAuth.RecoveryRequest{Username: username}
	res, err := authHandler.grpcClient.AuthClient.Recovery(context.Background(), &r)
	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}
	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) Recover(ctx *gin.Context) {

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

	res, err := authHandler.grpcClient.AuthClient.Recover(context.Background(), &rl)

	if err != nil {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusOK, &res)
}

func (authHandler *AuthHandler) Register(ctx *gin.Context) {

	registerDto := dto.RegisterDTO{}

	if err := ctx.BindJSON(&registerDto); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	v := validator.New()
	validators.UsernameValidator(ctx, v)
	validators.EmailValidator(ctx, v)
	validators.PasswordValidator(ctx, v)
	validators.NameValidator(ctx, v)
	validators.NumberValidator(ctx, v)
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

	errProfile := authHandler.registerProfile(userID, registerDto)
	if errProfile != nil {
		ctx.AbortWithError(http.StatusBadGateway, errProfile)
		return
	}

	errConnection := authHandler.registerConnection(userID, registerDto.IsPrivate)
	if errProfile != nil {
		ctx.AbortWithError(http.StatusBadGateway, errConnection)
		return
	}

	fmt.Println("successfully registered new user with ID:", userID)

	responsDTO := dto.RegisterResponsDTO{Id: userID, Username: registerDto.Username}

	ctx.JSON(http.StatusCreated, &responsDTO)
}

func (authHandler *AuthHandler) SendMailForMagicLinkRegistration(ctx *gin.Context) {
	authService := authHandler.grpcClient.AuthClient
	b := pbAuth.EmailForPasswordlessLoginRequest{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	response, _ := authService.SendEmailForPasswordlessLogin(ctx, &b)
	ctx.JSON(http.StatusOK, &response)
}

func (authHandler *AuthHandler) registerAuth(registerDTO dto.RegisterDTO) (string, error) {
	authS := authHandler.grpcClient.AuthClient
	response, err := authS.Register(context.TODO(), &pbAuth.RegisterRequest{Username: registerDTO.Username, Password: registerDTO.Password, Email: registerDTO.Email})
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

func (authHandler *AuthHandler) MagicLinkLogin(ctx *gin.Context) {
	authService := authHandler.grpcClient.AuthClient
	passwordlessMessage := pbAuth.PasswordlessLoginRequest{}
	if err := ctx.BindJSON(&passwordlessMessage); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	login, err := authService.PasswordlessLogin(ctx, &passwordlessMessage)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &login)
		return
	}
	ctx.JSON(http.StatusOK, &login)
}

func (authHandler *AuthHandler) GenerateApiToken(ctx *gin.Context) {
	userId := ctx.Param("userId")
	authService := authHandler.grpcClient.AuthClient
	res, err := authService.GenerateApiToken(ctx, &pbAuth.ApiTokenRequest{UserId: userId})
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
	ctx.JSON(http.StatusOK, "Everything is OK")
}
