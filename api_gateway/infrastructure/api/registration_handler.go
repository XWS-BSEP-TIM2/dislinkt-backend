package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/DTO"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/services"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

type RegistrationHandler struct {
	authClientAddress        string
	profileClientAddress     string
	connectionsClientAddress string
}

func NewRegistrationHandler(authClientAddress, profileClientAddress, connectionsClientAddress string) Handler {
	return &RegistrationHandler{
		authClientAddress:        authClientAddress,
		profileClientAddress:     profileClientAddress,
		connectionsClientAddress: connectionsClientAddress,
	}
}

func (handler *RegistrationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/api/register", handler.Register)
	if err != nil {
		panic(err)
	}
}

func (handler *RegistrationHandler) Register(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	var registerDTO DTO.RegisterDTO

	err := json.NewDecoder(r.Body).Decode(&registerDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Registration: ", registerDTO)

	userID, errAuth := handler.RegisterAuth(registerDTO)
	if errAuth != nil {
		http.Error(w, errAuth.Error(), http.StatusBadRequest)
		return
	}

	errProfile := handler.RegisterProfile(userID, registerDTO)
	if errProfile != nil {
		http.Error(w, errProfile.Error(), http.StatusBadRequest)
		return
	}

	errConnection := handler.RegisterConnection(userID, registerDTO.IsPrivate)
	if errConnection != nil {
		http.Error(w, errConnection.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("successfully registered new user with ID:", userID)

	response, errRes := json.Marshal(&DTO.RegisterResponsDTO{Id: userID, Username: registerDTO.Username})
	if errRes != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (handler *RegistrationHandler) RegisterAuth(registerDTO DTO.RegisterDTO) (string, error) {
	authS := services.NewAuthClient(handler.authClientAddress)
	response, err := authS.Register(context.TODO(), &pbAuth.RegisterRequest{Data: &pbAuth.RegisterRequestData{Username: registerDTO.Username, Password: registerDTO.Password}}) //TODO: izbaciti neiskrostene atribute
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(response)
	return response.UserID, err
}

func (handler *RegistrationHandler) RegisterProfile(userID string, registerDTO DTO.RegisterDTO) error {
	profileService := services.NewProfileClient(handler.profileClientAddress)
	response, err := profileService.CreateProfile(context.TODO(), &pbProfile.CreateProfileRequest{Profile: registerDTO.ToProto(userID)})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(response)
	return nil
}

func (handler *RegistrationHandler) RegisterConnection(userID string, IsPrivate bool) error {
	connectionService := services.NewConnectionClient(handler.connectionsClientAddress)
	registrationResult, err := connectionService.Register(context.TODO(), &pbConnection.RegisterRequest{User: &pbConnection.User{UserID: userID, IsPublic: IsPrivate}})
	fmt.Println(registrationResult)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
