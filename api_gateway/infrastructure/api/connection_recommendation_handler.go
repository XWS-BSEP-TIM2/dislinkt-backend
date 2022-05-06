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

type ConnectionRecommendationHandler struct {
	authClientAddress        string
	profileClientAddress     string
	connectionsClientAddress string
}

func NewConnectionRecommendationHandler(authClientAddress, profileClientAddress, connectionsClientAddress string) Handler {
	return &ConnectionRecommendationHandler{
		authClientAddress:        authClientAddress,
		profileClientAddress:     profileClientAddress,
		connectionsClientAddress: connectionsClientAddress,
	}
}

func (handler *ConnectionRecommendationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/api/connection/recommendation", handler.Recommendation)
	if err != nil {
		panic(err)
	}
}

func (handler *ConnectionRecommendationHandler) Recommendation(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	token, err1 := apiUtils.GetJWTfromHTTPReq(r)
	if err1 != nil {
		fmt.Println(err1.Error())
		http.Error(w, err1.Error(), http.StatusUnauthorized)
		return
	}

	authS := services.NewAuthClient(handler.authClientAddress)
	dataFromToken, err2 := authS.ExtractDataFromToken(context.TODO(), &pbAuth.ExtractDataFromTokenRequest{Token: token})
	if err2 != nil {
		fmt.Println(err2.Error())
		http.Error(w, err2.Error(), http.StatusUnauthorized)
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

	profileService := services.NewProfileClient(handler.profileClientAddress)

	var recommendationDTO []*DTO.ConnectionDTO

	for _, v := range recommendations.Users {
		profile, err := profileService.Get(ctx, &pbProfile.GetRequest{Id: v.UserID})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		recommendationDTO = append(recommendationDTO, DTO.MapConnectionDTO(profile.Profile))
	}

	response, errRes := json.Marshal(recommendationDTO)
	if errRes != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
