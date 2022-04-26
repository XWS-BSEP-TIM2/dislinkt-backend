package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	//cfg "github.com/tamararankovic/microservices_demo/api_gateway/startup/config"
	cfg "api_gateway/startup/config"

	postGw "github.com/tamararankovic/microservices_demo/common/proto/post_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	config *cfg.Config
	mux    *runtime.ServeMux
}

func NewServer(config *cfg.Config) *Server {
	server := &Server{
		config: config,
		mux:    runtime.NewServeMux(),
	}
	server.initHandlers()
	server.initCustomHandlers()
	return server
}

func (server *Server) initHandlers() {

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	PostEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)
	err := postGw.RegisterPostServiceHandlerFromEndpoint(context.TODO(), server.mux, PostEndpoint, opts)
	if err != nil {
		panic(err)
	}

}

func (server *Server) initCustomHandlers() {

	/*
		catalogueEmdpoint := fmt.Sprintf("%s:%s", server.config.CatalogueHost, server.config.CataloguePort)
		orderingEmdpoint := fmt.Sprintf("%s:%s", server.config.OrderingHost, server.config.OrderingPort)
		shippingEmdpoint := fmt.Sprintf("%s:%s", server.config.ShippingHost, server.config.ShippingPort)
	*/
	/*
		orderingHandler := api.NewPostHandler("localhost:8000")
		orderingHandler.Init(server.mux)
	*/
}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}
