package startup

import (
	"fmt"
	postGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	postStore := server.initPostStore(mongoClient)
	postService := server.initPostService(postStore)
	postHandler := server.initPostHandler(postService)
	server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initPostStore(client *mongo.Client) domain.PostStore {
	return persistence.NewPostMongoDBStore(client)
}

func (server *Server) initPostService(store domain.PostStore) *application.PostService {
	authServiceAddress := fmt.Sprintf("%s:%s", server.config.AuthServiceHost, server.config.AuthServicePort)
	connectionServiceAddress := fmt.Sprintf("%s:%s", server.config.ConnectionServiceHost, server.config.ConnectionServicePort)
	profileServiceAddress := fmt.Sprintf("%s:%s", server.config.ProfileServiceHost, server.config.ProfileServicePort)

	return application.NewPostService(store, authServiceAddress, connectionServiceAddress, profileServiceAddress)
}

func (server *Server) initPostHandler(service *application.PostService) *api.PostHandler {
	return api.NewPostHandler(service)
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(
	//	grpc_middleware.ChainUnaryServer(
	//		interceptors.TokenAuthInterceptor,
	//	),
	//),
	)
	postGw.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
