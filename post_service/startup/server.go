package startup

import (
	"crypto/tls"
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	postGw "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/post_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup/config"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/post_service/startup/data"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	trace, _ := tracer.Init("post_service")
	opentracing.SetGlobalTracer(trace)

	mongoClient := server.initMongoClient()
	loggingService := server.initLoggingService()
	postStore := server.initPostStore(mongoClient, loggingService)
	postService := server.initPostService(postStore)
	postHandler := server.initPostHandler(postService)

	fmt.Println("Post service started.")
	server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initPostStore(client *mongo.Client, loggingService loggingS.LoggingServiceClient) domain.PostStore {
	store := persistence.NewPostMongoDBStore(client, loggingService)
	data.InitializePostStore(store)
	return store
}

func (server *Server) initPostService(store domain.PostStore) *application.PostService {
	authServiceAddress := fmt.Sprintf("%s:%s", server.config.AuthServiceHost, server.config.AuthServicePort)
	connectionServiceAddress := fmt.Sprintf("%s:%s", server.config.ConnectionServiceHost, server.config.ConnectionServicePort)
	profileServiceAddress := fmt.Sprintf("%s:%s", server.config.ProfileServiceHost, server.config.ProfileServicePort)
	notificationServiceAddress := fmt.Sprintf("%s:%s", server.config.NotificationServiceHost, server.config.NotificationServicePort)
	loggingServiceAddress := fmt.Sprintf("%s:%s", server.config.LoggingServiceHost, server.config.LoggingServicePort)

	return application.NewPostService(
		store,
		authServiceAddress,
		connectionServiceAddress,
		profileServiceAddress,
		notificationServiceAddress,
		loggingServiceAddress)
}

func (server *Server) initPostHandler(service *application.PostService) *api.PostHandler {
	return api.NewPostHandler(service)
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	creds, err := credentials.NewServerTLSFromFile("./certificates/post_service.crt", "./certificates/post_service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	postGw.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initLoggingService() loggingS.LoggingServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.LoggingServiceHost, server.config.LoggingServicePort)
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Logging service: %v", err)
	}
	return loggingS.NewLoggingServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}
