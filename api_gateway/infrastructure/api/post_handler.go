package api

// legacy - should be removed

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/services"
	post "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcMetadata "google.golang.org/grpc/metadata"
	"net/http"
)

type PostHandler struct {
	postClientAddress string
}

func NewPostHandler(orderingClientAddress string) Handler {
	return &PostHandler{
		postClientAddress: orderingClientAddress,
	}
}

func (handler *PostHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/post-gateway/{postId}", handler.GetDetails)
	if err != nil {
		panic(err)
	}
}

func (handler *PostHandler) GetDetails(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["postId"]
	fmt.Println("ID: ", id)
	postClient := services.NewPostClient(handler.postClientAddress)
	var ctx = context.TODO()
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", r.Header.Get("authorization"))

	post, err := postClient.GetPost(ctx, &post.GetPostRequest{PostId: id})
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported) //208
		fmt.Println("POST: ", post)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)

}
