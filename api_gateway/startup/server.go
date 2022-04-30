package startup

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/api"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	cfg "github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"

	authGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	connectionGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	postGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	profileGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
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

	AuthEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	err = authGw.RegisterAuthServiceHandlerFromEndpoint(context.TODO(), server.mux, AuthEndpoint, opts)
	if err != nil {
		panic(err)
	}
	ProfileEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileHost, server.config.ProfilePort)
	err = profileGw.RegisterProfileServiceHandlerFromEndpoint(context.TODO(), server.mux, ProfileEndpoint, opts)
	if err != nil {
		panic(err)
	}

	ConnectionEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	err = connectionGw.RegisterConnectionServiceHandlerFromEndpoint(context.TODO(), server.mux, ConnectionEndpoint, opts)
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

	orderingHandler := api.NewPostHandler("localhost:8080")
	orderingHandler.Init(server.mux)

}

func (server *Server) Start() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.mux))
}
