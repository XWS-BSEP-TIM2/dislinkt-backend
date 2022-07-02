package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/converter"
	pbLogg "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	events "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/create_order"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"google.golang.org/grpc/peer"
	"net/http"
	"time"

	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/utils"
	pb "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	userService              *application.AuthService
	Jwt                      utils.JwtWrapper
	passwordlessService      *application.PasswordlessTokenService
	apiTokenService          *application.ApiTokenService
	LoggingService           pbLogg.LoggingServiceClient
	RegisterUserOrchestrator *application.RegisterUserOrchestrator
}

func NewAuthHandler(service *application.AuthService, passwordlessServices *application.PasswordlessTokenService, apiTokenService *application.ApiTokenService, loggingService pbLogg.LoggingServiceClient, orchestrator *application.RegisterUserOrchestrator) *AuthHandler {
	return &AuthHandler{
		userService:              service,
		passwordlessService:      passwordlessServices,
		apiTokenService:          apiTokenService,
		LoggingService:           loggingService,
		RegisterUserOrchestrator: orchestrator,
	}
}

func (handler *AuthHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Register")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	var user domain.User
	user1, _ := handler.userService.GetByUsername(ctx2, request.Username)
	if user1 != nil {
		handler.logg(ctx2, "ERROR", "Register", request.Username, "Username is not unique")
		return &pb.RegisterResponse{
			Status: http.StatusUnprocessableEntity,
			Error:  "Username is not unique",
			UserID: "",
		}, nil
	}
	user.Username = request.GetUsername()
	user.Password = utils.HashPassword(request.GetPassword())
	fmt.Println(user.Password)
	user.Email = request.GetEmail()

	token, err := utils.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	user.VerificationCode = token
	user.VerificationCodeTime = time.Now()

	userID, err := handler.userService.Create(ctx2, &user) //userID
	if err != nil {
		return &pb.RegisterResponse{
			Status: http.StatusUnauthorized,
			UserID: "",
		}, err
	}

	user.Id = converter.GetObjectId(userID)
	handler.RegisterUserOrchestrator.Start(events.UserDetails{Id: user.Id.Hex(), Birthday: request.BirthDate.AsTime(), Surname: request.Surname, Username: request.Username, Email: request.Email, Gender: request.Gender, PhoneNumber: request.PhoneNumber, IsPrivate: request.IsPrivate, Name: request.Name})

	errSendVerification := handler.userService.SendVerification(ctx2, &user)
	if errSendVerification != nil {
		fmt.Println("Error:", errSendVerification.Error())
	}

	handler.logg(ctx2, "SUCCESS", "Register", request.Username, "Successfully register new user")
	return &pb.RegisterResponse{
		Status: http.StatusCreated,
		UserID: userID,
	}, nil

}

func (handler *AuthHandler) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Verify")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	ret, err := handler.userService.Verify(ctx2, req.Username, req.Code)
	return ret, err
}

func (handler *AuthHandler) ResendVerify(ctx context.Context, req *pb.ResendVerifyRequest) (*pb.ResendVerifyResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ResendVerify")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.userService.ResendVerify(ctx2, req.Username)
}

func (handler *AuthHandler) Recovery(ctx context.Context, req *pb.RecoveryRequest) (*pb.RecoveryResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Recovery")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	return handler.userService.Recovery(ctx2, req.Username)
}

func (handler *AuthHandler) Recover(ctx context.Context, req *pb.RecoveryRequestLogin) (*pb.LoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Recover")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	response, err := handler.userService.Recover(ctx2, req)
	if err != nil {
		return response, err
	}
	if response.Error != "" {
		return response, nil
	}
	return handler.Login(ctx2, &pb.LoginRequest{Username: req.Username, Password: req.NewPassword})
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ChangePassword")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	ret, err := handler.userService.ChangePassword(ctx2, req)
	if ret.Status == 200 {
		handler.logg(ctx2, "SUCCESS", "ChangePassword", req.Username, ret.Msg)
	} else {
		handler.logg(ctx2, "ERROR", "ChangePassword", req.Username, ret.Msg)
	}
	return ret, err
}

func (handler *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Login")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	user, err := handler.userService.GetByUsername(ctx2, req.Username)
	if err != nil {
		handler.logg(ctx2, "WARNING", "Login", "", "Username or password is incorrect")
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Username or password is incorrect",
		}, nil
	}

	if user.Locked {
		handler.logg(ctx2, "WARNING", "Login", user.Id.Hex(), "Acc is locked")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  user.LockReason,
		}, nil
	}

	if !user.Verified {
		handler.logg(ctx2, "WARNING", "Login", user.Id.Hex(), "Your Acc is not verified")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  "Your Acc is not verified",
		}, nil
	}

	if user.NumOfErrTryLogin == 5 && !user.LastErrTryLoginTime.Add(1*time.Hour).Before(time.Now()) {
		handler.logg(ctx2, "WARNING", "Login", user.Id.Hex(), fmt.Sprint(user.NumOfErrTryLogin)+" failed login attempts")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(1*time.Hour).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 4 && !user.LastErrTryLoginTime.Add(15*time.Minute).Before(time.Now()) {
		handler.logg(ctx2, "WARNING", "Login", user.Id.Hex(), fmt.Sprint(user.NumOfErrTryLogin)+" failed login attempts")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(15*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 3 && !user.LastErrTryLoginTime.Add(3*time.Minute).Before(time.Now()) {
		handler.logg(ctx2, "WARNING", "Login", user.Id.Hex(), fmt.Sprint(user.NumOfErrTryLogin)+" failed login attempts")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(3*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		handler.logg(ctx2, "ERROR", "Login", user.Id.Hex(), "Username or password is incorrect")
		user.NumOfErrTryLogin += 1
		user.LastErrTryLoginTime = time.Now()
		if user.NumOfErrTryLogin >= 6 {
			user.Locked = true
			user.LockReason = "your account is locked, due to many incorrect login attempts"
			handler.logg(ctx2, "INFO", "Login", user.Id.Hex(), "account is locked, due to many incorrect login attempts")
		}
		handler.userService.Update(ctx2, user)
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Username or password is incorrect",
		}, nil
	}

	if user.IsTFAEnabled {
		return &pb.LoginResponse{
			Status:    http.StatusOK,
			Role:      domain.ConvertRoleToString(user.Role),
			Username:  user.Username,
			UserID:    user.Id.Hex(),
			TwoFactor: true,
		}, nil
	}

	token, _ := handler.Jwt.GenerateToken(user)

	user.NumOfErrTryLogin = 0
	handler.userService.Update(ctx2, user)

	handler.logg(ctx2, "SUCCESS", "Login", user.Id.Hex(), "successfully login")
	return &pb.LoginResponse{
		Status:    http.StatusOK,
		Token:     token,
		Role:      domain.ConvertRoleToString(user.Role),
		Username:  user.Username,
		UserID:    user.Id.Hex(),
		TwoFactor: false,
	}, nil
}

func (handler *AuthHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Validate")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	claims, err := handler.Jwt.ValidateToken(req.Token)

	if err != nil {
		handler.logg(ctx2, "ERROR", "Validate", "", "Invalid JWT")
		return &pb.ValidateResponse{
			Status: http.StatusBadRequest,
			Error:  err.Error(),
		}, nil
	}

	user, err := handler.userService.Get(ctx2, getObjectId(claims.Id))
	if err != nil {
		return &pb.ValidateResponse{
			Status: http.StatusNotFound,
			Error:  "User not found",
		}, nil
	}

	return &pb.ValidateResponse{
		Status: http.StatusOK,
		UserId: user.Id.Hex(),
	}, nil
}

func (handler *AuthHandler) ExtractDataFromToken(ctx context.Context, req *pb.ExtractDataFromTokenRequest) (*pb.ExtractDataFromTokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ExtractDataFromToken")
	defer span.Finish()

	claims, err := handler.Jwt.ValidateToken(req.Token)

	if err != nil {
		return &pb.ExtractDataFromTokenResponse{
			Id:       "",
			Username: "",
			Role:     "",
		}, err
	}

	return &pb.ExtractDataFromTokenResponse{
		Id:        claims.Id,
		Username:  claims.Username,
		Role:      claims.Role,
		TokenType: claims.TokenType,
		ApiCode:   claims.ApiCode,
	}, nil

}

func (handler *AuthHandler) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.LoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "PasswordlessLogin")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	token, err := handler.passwordlessService.GetByTokenCode(ctx2, request.TokenCode)
	if err != nil {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", "", "Token does not exist")
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Token does not exist",
		}, nil
	}
	if token.CreationDate.Add(15 * time.Minute).Before(time.Now()) {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", "", "Token is expired")
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Token is expired",
		}, err
	}

	handler.passwordlessService.Delete(ctx2, token.Id)

	user, err := handler.userService.Get(ctx2, token.UserId)
	if err != nil {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", user.Id.Hex(), "Username or password is incorrect")
		return &pb.LoginResponse{
			Status: http.StatusNotFound,
			Error:  "Username or password is incorrect",
		}, nil
	}

	if user.Locked {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", user.Id.Hex(), "Acc is locked")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  user.LockReason,
		}, nil
	}

	if !user.Verified {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", user.Id.Hex(), "Your Acc is not verified")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  "Your Acc is not verified",
		}, nil
	}

	if user.NumOfErrTryLogin == 5 && !user.LastErrTryLoginTime.Add(1*time.Hour).Before(time.Now()) {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", user.Id.Hex(), fmt.Sprint(user.NumOfErrTryLogin)+" failed login attempts")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(1*time.Hour).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 4 && !user.LastErrTryLoginTime.Add(15*time.Minute).Before(time.Now()) {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", user.Id.Hex(), fmt.Sprint(user.NumOfErrTryLogin)+" failed login attempts")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(15*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	} else if user.NumOfErrTryLogin == 3 && !user.LastErrTryLoginTime.Add(3*time.Minute).Before(time.Now()) {
		handler.logg(ctx2, "ERROR", "PasswordlessLogin", user.Id.Hex(), fmt.Sprint(user.NumOfErrTryLogin)+" failed login attempts")
		return &pb.LoginResponse{
			Status: http.StatusForbidden,
			Error:  fmt.Sprint(user.NumOfErrTryLogin) + " failed login attempts, you will be able to login after " + fmt.Sprintf("%f", user.LastErrTryLoginTime.Add(3*time.Minute).Sub(time.Now()).Minutes()) + " minutes",
		}, nil
	}

	tokenJwt, _ := handler.Jwt.GenerateToken(user)

	handler.logg(ctx2, "SUCCESS", "PasswordlessLogin", user.Id.Hex(), "Successfully login")
	return &pb.LoginResponse{
		Status:   http.StatusOK,
		Token:    tokenJwt,
		Role:     domain.ConvertRoleToString(user.Role),
		Username: user.Username,
		UserID:   user.Id.Hex(),
	}, nil
}

func (handler *AuthHandler) SendEmailForPasswordlessLogin(ctx context.Context, request *pb.EmailForPasswordlessLoginRequest) (*pb.SendEmailForPasswordLoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SendEmailForPasswordlessLogin")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	user, err := handler.userService.GetByEmail(ctx2, request.Email)
	if err != nil || user == nil {
		return &pb.SendEmailForPasswordLoginResponse{
			Error: "Email does not exist",
		}, err
	}
	tokenCode, _ := utils.GenerateRandomStringURLSafe(30)
	token := domain.PasswordlessToken{
		TokenCode:    tokenCode,
		UserId:       user.Id,
		CreationDate: time.Now(),
	}
	handler.passwordlessService.Create(ctx2, &token)
	handler.passwordlessService.SendMagicLink(ctx2, user, tokenCode)
	handler.logg(ctx2, "SUCCESS", "SendEmailForPasswordlessLogin", user.Id.Hex(), "Successfully send magic link")
	return &pb.SendEmailForPasswordLoginResponse{
		Error: "",
	}, err
}

func getObjectId(id string) primitive.ObjectID {
	if objectId, err := primitive.ObjectIDFromHex(id); err == nil {
		return objectId
	}
	return primitive.NewObjectID()
}

func (handler *AuthHandler) GenerateApiToken(ctx context.Context, request *pb.ApiTokenRequest) (*pb.ApiTokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GenerateApiToken")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	user, _ := handler.userService.Get(ctx2, getObjectId(request.UserId))
	if user == nil {
		handler.logg(ctx2, "ERROR", "GenerateApiToken", user.Id.Hex(), "User ID does not exist")
		error := pb.ErrorResponse{ErrorCode: 500, Message: "User ID does not exist"}
		return &pb.ApiTokenResponse{Error: &error}, nil
	}
	tokenCode, err := handler.apiTokenService.Create(ctx2, getObjectId(request.UserId))
	if err != nil {
		handler.logg(ctx2, "ERROR", "GenerateApiToken", user.Id.Hex(), "Unable to create api token")
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Unable to create api token"}
		return &pb.ApiTokenResponse{Error: &error}, nil
	}

	handler.logg(ctx2, "SUCCESS", "GenerateApiToken", user.Id.Hex(), "Successfully generated api token")
	return &pb.ApiTokenResponse{TokenCode: tokenCode, Error: nil}, nil
}

func (handler *AuthHandler) ValidateApiToken(ctx context.Context, request *pb.ValidateApiTokenRequest) (*pb.ValidateApiTokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ValidateApiToken")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	token, err := handler.apiTokenService.GetByTokenCode(ctx2, request.TokenCode)
	if err != nil {
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Error while searching token"}
		return &pb.ValidateApiTokenResponse{Error: &error}, nil
	}
	if token == nil {
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Invalid token"}
		return &pb.ValidateApiTokenResponse{Error: &error}, nil
	}
	return &pb.ValidateApiTokenResponse{Error: nil}, nil
}

func (handler *AuthHandler) GetApiToken(ctx context.Context, request *pb.GetApiTokenRequest) (*pb.GetApiTokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetApiToken")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	token, err := handler.apiTokenService.GetByTokenCode(ctx2, request.TokenCode)
	if err != nil {
		return nil, err
	}
	tokenProto := pb.ApiToken{TokenCode: token.ApiCode, UserId: token.UserId.Hex()}
	return &pb.GetApiTokenResponse{Token: &tokenProto}, nil
}

func (handler *AuthHandler) GenerateQr2TF(ctx context.Context, request *pb.UserIdRequest) (*pb.TFAResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GenerateQr2TF")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	qrCode, err := handler.userService.GenerateQR2FA(ctx2, converter.GetObjectId(request.UserId))
	if err != nil {
		handler.logg(ctx2, "ERROR", "GenerateQr2TF", request.UserId, "Unable to generate qr code")
		error := pb.ErrorResponse{ErrorCode: 500, Message: "Unable to generate qr code"}
		return &pb.TFAResponse{Error: &error}, nil
	}
	handler.logg(ctx2, "SUCCESS", "GenerateQr2TF", request.UserId, "Successfully generated qr code")
	return &pb.TFAResponse{QrCode: qrCode}, nil
}

func (handler *AuthHandler) Verify2FactorCode(ctx context.Context, request *pb.TFARequest) (*pb.LoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Verify2FactorCode")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	err := handler.userService.Verify2fa(ctx2, converter.GetObjectId(request.UserId), request.Code)
	if err != nil {
		handler.logg(ctx2, "ERROR", "Verify2FactorCode", request.UserId, "Wrong code")
		return &pb.LoginResponse{Status: 401, Error: "Wrong code"}, nil
	}

	user, err := handler.userService.Get(ctx2, converter.GetObjectId(request.UserId))
	if err != nil {
		handler.logg(ctx2, "ERROR", "Verify2FactorCode", request.UserId, "Wrong code")
		return &pb.LoginResponse{Status: 401, Error: "Wrong code"}, nil
	}

	token, _ := handler.Jwt.GenerateToken(user)

	user.NumOfErrTryLogin = 0
	handler.userService.Update(ctx2, user)

	handler.logg(ctx2, "SUCCESS", "Verify2FactorCode", request.UserId, "Successfully verified two factor auth code")
	return &pb.LoginResponse{
		Status:    http.StatusOK,
		Token:     token,
		Role:      domain.ConvertRoleToString(user.Role),
		Username:  user.Username,
		UserID:    user.Id.Hex(),
		TwoFactor: false,
	}, nil

}

func (handler *AuthHandler) EditData(ctx context.Context, request *pb.EditDataRequest) (*pb.EditDataResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "EditData")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	userCredentials, _ := handler.userService.Get(ctx2, converter.GetObjectId(request.UserId))
	userCredentials.Username = request.Username
	userCredentials.IsTFAEnabled = request.IsTwoFactor
	userCredentials.Email = request.Email

	handler.userService.Update(ctx2, userCredentials)
	return &pb.EditDataResponse{}, nil
}

func (handler *AuthHandler) logg(ctx context.Context, logType, serviceFunctionName, userID, description string) {
	span := tracer.StartSpanFromContext(ctx, "logg")
	defer span.Finish()
	ctx2 := tracer.ContextWithSpan(context.Background(), span)

	ipAddress := ""
	p, ok := peer.FromContext(ctx)
	if ok {
		ipAddress = p.Addr.String()
	}
	if logType == "ERROR" {
		handler.LoggingService.LoggError(ctx2, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "SUCCESS" {
		handler.LoggingService.LoggSuccess(ctx2, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "WARNING" {
		handler.LoggingService.LoggWarning(ctx2, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	} else if logType == "INFO" {
		handler.LoggingService.LoggInfo(ctx2, &pbLogg.LogRequest{ServiceName: "AUTH_SERVICE", ServiceFunctionName: serviceFunctionName, UserID: userID, IpAddress: ipAddress, Description: description})
	}
}
