package startup

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/infrastructure/api"
	cfg "github.com/XWS-BSEP-TIM2/dislinkt-backend/api_gateway/startup/config"
	authGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	connectionGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	postGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	profileGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

	authEndpoint := fmt.Sprintf("%s:%s", server.config.AuthHost, server.config.AuthPort)
	profileEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileHost, server.config.ProfilePort)
	connectionsEndpoint := fmt.Sprintf("%s:%s", server.config.ConnectionHost, server.config.ConnectionPort)
	postEndpoint := fmt.Sprintf("%s:%s", server.config.PostHost, server.config.PostPort)

	postHandler := api.NewPostHandler(postEndpoint)
	postHandler.Init(server.mux)

	profileHandler := api.NewPostHandler(profileEndpoint)
	profileHandler.Init(server.mux)

	registerHandler := api.NewRegistrationHandler(authEndpoint, profileEndpoint, connectionsEndpoint)
	registerHandler.Init(server.mux)

	connectionRecommendationHandler := api.NewConnectionHandler(authEndpoint, profileEndpoint, connectionsEndpoint)
	connectionRecommendationHandler.Init(server.mux)

}

func (server *Server) Start() {

	ch := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://localhost:4200", "https://localhost:4200/**", "http://localhost:4200", "http://localhost:4200/**"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "Authorization", "Access-Control-Allow-Origin", "*"}),
	)

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", server.config.Port), server.config.CertificatePath, server.config.CertificatePrivateKeyPath, ch(server.mux)))
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), ch(server.mux)))
}
