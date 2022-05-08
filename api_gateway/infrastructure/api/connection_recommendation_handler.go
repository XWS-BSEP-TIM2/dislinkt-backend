package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/DTO"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/apiUtils"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/services"
	pbAuth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	pbConnection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbProfile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcMetadata "google.golang.org/grpc/metadata"
	"net/http"
)

type ConnectionHandler struct {
	authClientAddress        string
	profileClientAddress     string
	connectionsClientAddress string
}

func NewConnectionHandler(authClientAddress, profileClientAddress, connectionsClientAddress string) Handler {
	return &ConnectionHandler{
		authClientAddress:        authClientAddress,
		profileClientAddress:     profileClientAddress,
		connectionsClientAddress: connectionsClientAddress,
	}
}

func (handler *ConnectionHandler) Init(mux *runtime.ServeMux) {
	err1 := mux.HandlePath("GET", "/api/connection/recommendation", handler.Recommendation)
	if err1 != nil {
		panic(err1)
	}

	err2 := mux.HandlePath("GET", "/api/connection/friends/{userID}", handler.GetFriends)
	if err2 != nil {
		panic(err2)
	}

	err3 := mux.HandlePath("GET", "/api/connection/block", handler.GetBlocks)
	if err3 != nil {
		panic(err3)
	}

	err4 := mux.HandlePath("GET", "/api/connection/friends-request", handler.GetFriendsRequest)
	if err4 != nil {
		panic(err4)
	}
}

func (handler *ConnectionHandler) GetBlocks(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	token, dataFromToken, err := handler.GetTokenAndData(w, r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := grpcMetadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+token)

	connectionService := services.NewConnectionClient(handler.connectionsClientAddress)
	recommendations, err3 := connectionService.GetBlockeds(ctx, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err3 != nil {
		fmt.Println(err3.Error())
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx, w, r, recommendations.Users)

	handler.MarshalData(w, r, recommendationDTO)
}

func (handler *ConnectionHandler) GetFriendsRequest(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	token, dataFromToken, err := handler.GetTokenAndData(w, r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := grpcMetadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+token)

	connectionService := services.NewConnectionClient(handler.connectionsClientAddress)
	recommendations, err3 := connectionService.GetFriendRequests(ctx, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err3 != nil {
		fmt.Println(err3.Error())
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx, w, r, recommendations.Users)

	handler.MarshalData(w, r, recommendationDTO)
}

func (handler *ConnectionHandler) Recommendation(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	token, dataFromToken, err := handler.GetTokenAndData(w, r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := grpcMetadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+token)

	connectionService := services.NewConnectionClient(handler.connectionsClientAddress)
	recommendations, err3 := connectionService.GetRecommendation(ctx, &pbConnection.GetRequest{UserID: dataFromToken.Id})
	if err3 != nil {
		fmt.Println(err3.Error())
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	recommendationDTO := handler.GenerateProfileDTO(ctx, w, r, recommendations.Users)

	handler.MarshalData(w, r, recommendationDTO)
}

func (handler *ConnectionHandler) GetFriends(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	token, _, err := handler.GetTokenAndData(w, r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	targetUserID := pathParams["userID"]

	ctx := grpcMetadata.AppendToOutgoingContext(r.Context(), "authorization", "Bearer "+token)

	connectionService := services.NewConnectionClient(handler.connectionsClientAddress)
	friends, err3 := connectionService.GetFriends(ctx, &pbConnection.GetRequest{UserID: targetUserID})
	if err3 != nil {
		fmt.Println(err3.Error())
		http.Error(w, err3.Error(), http.StatusInternalServerError)
		return
	}

	friendsDTO := handler.GenerateProfileDTO(ctx, w, r, friends.Users)

	handler.MarshalData(w, r, friendsDTO)
}

// HELP fun

func (handler *ConnectionHandler) GetTokenAndData(w http.ResponseWriter, r *http.Request) (string, *pbAuth.ExtractDataFromTokenResponse, error) {
	token, err1 := apiUtils.GetJWTfromHTTPReq(r)
	if err1 != nil {
		fmt.Println(err1.Error())
		http.Error(w, err1.Error(), http.StatusUnauthorized)
		return "", nil, err1
	}

	authS := services.NewAuthClient(handler.authClientAddress)
	dataFromToken, err2 := authS.ExtractDataFromToken(context.TODO(), &pbAuth.ExtractDataFromTokenRequest{Token: token})
	if err2 != nil {
		fmt.Println(err2.Error())
		http.Error(w, err2.Error(), http.StatusUnauthorized)
		return token, nil, err2
	}

	return token, dataFromToken, nil
}

func (handler *ConnectionHandler) MarshalData(w http.ResponseWriter, r *http.Request, connections []*DTO.ConnectionDTO) {
	response, errRes := json.Marshal(connections)
	if errRes != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ConnectionHandler) GenerateProfileDTO(ctx context.Context, w http.ResponseWriter, r *http.Request, connUsers []*pbConnection.User) []*DTO.ConnectionDTO {
	profileService := services.NewProfileClient(handler.profileClientAddress)

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
