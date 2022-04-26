package api

import (
	"api_gateway/infrastructure/services"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	post "github.com/tamararankovic/microservices_demo/common/proto/post_service"
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
	err := mux.HandlePath("GET", "/post/{postId}", handler.GetDetails)
	if err != nil {
		panic(err)
	}
}

func (handler *PostHandler) GetDetails(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["postId"]
	fmt.Println("ID: ", id)

	postClient := services.NewPostClient(handler.postClientAddress)
	post, err := postClient.Get(context.TODO(), &post.GetRequest{Id: id})
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
