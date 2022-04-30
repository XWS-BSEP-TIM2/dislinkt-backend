package startup

import (
	"fmt"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/interceptors"
	post "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup/config"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
	productStore := server.initProductStore(mongoClient)

	productService := server.initProductService(productStore)

	productHandler := server.initProductHandler(productService)

	server.startGrpcServer(productHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initProductStore(client *mongo.Client) domain.PostStore {
	store := persistence.NewProductMongoDBStore(client)
	store.DeleteAll()
	for _, product := range products {
		err := store.Insert(product)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initProductService(store domain.PostStore) *application.PostService {
	return application.NewProductService(store)
}

func (server *Server) initProductHandler(service *application.PostService) *api.PostHandler {
	return api.NewProductHandler(service)
}

func (server *Server) startGrpcServer(productHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				interceptors.TokenAuthInterceptor,
			),
		),
	)
	post.RegisterPostServiceServer(grpcServer, productHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
