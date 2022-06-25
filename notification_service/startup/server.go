package startup

import (
	"context"
	"crypto/tls"
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	messageS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/startup/config"
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

	mongoClient := server.initMongoClient()

	loggingService := server.initLoggingService()

	profileStore := server.initMessageStore(mongoClient)

	profileService := server.initMessageService(profileStore, loggingService)

	profileHandler := server.initMessageHandler(profileService)

	fmt.Println("Notification service started.")
	server.startGrpcServer(profileHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.NotificationDBHost, server.config.NotificationDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initMessageStore(client *mongo.Client) persistence.NotificationStore {
	store := persistence.NewNotificationMongoDbStore(client)

	store.DeleteAll(context.TODO())
	for _, chat := range chats {
		err := store.Insert(context.TODO(), chat)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
}

func (server *Server) initMessageService(store persistence.NotificationStore, loggingService loggingS.LoggingServiceClient) *application.NotificationService {
	return application.NewNotificationService(store, server.config, loggingService)
}

func (server *Server) initMessageHandler(service *application.NotificationService) *api.NotificationHandler {
	return api.NewMessageHandler(service)
}

func (server *Server) startGrpcServer(profileHandler *api.NotificationHandler) {
	creds, err := credentials.NewServerTLSFromFile("./certificates/message_service.crt", "./certificates/message_service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	messageS.RegisterMessageServiceServer(grpcServer, profileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initLoggingService() loggingS.LoggingServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.LoggingHost, server.config.LoggingPort)
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
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
