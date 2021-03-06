package startup

import (
	"context"
	"crypto/tls"
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	messageS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/startup/config"
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
	trace, _ := tracer.Init("message_service")
	opentracing.SetGlobalTracer(trace)

	mongoClient := server.initMongoClient()

	loggingService := server.initLoggingService()

	messageStore := server.initMessageStore(mongoClient)

	messageService := server.initMessageService(messageStore, loggingService)

	messageHandler := server.initMessageHandler(messageService)

	fmt.Println("Message service started.")
	server.startGrpcServer(messageHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.MessageDBHost, server.config.MessageDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initMessageStore(client *mongo.Client) persistence.MessageStore {
	store := persistence.NewMessageMongoDbStore(client)

	store.DeleteAll(context.TODO())
	for _, chat := range chats {
		_, err := store.Insert(context.TODO(), chat)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
}

func (server *Server) initMessageService(store persistence.MessageStore, loggingService loggingS.LoggingServiceClient) *application.MessageService {
	return application.NewMessageService(store, server.config, loggingService)
}

func (server *Server) initMessageHandler(service *application.MessageService) *api.MessageHandler {
	return api.NewMessageHandler(service)
}

func (server *Server) startGrpcServer(profileHandler *api.MessageHandler) {
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
