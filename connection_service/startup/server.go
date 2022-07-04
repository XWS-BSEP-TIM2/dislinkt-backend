package startup

import (
	"crypto/tls"
	"fmt"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging/nats"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/opentracing/opentracing-go"
	"log"
	"net"

	connection "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/connection_service"
	pbLogg "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	pbMessage "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/connection_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "connection_service"
)

func (server *Server) Start() {
	trace, _ := tracer.Init("connection_service")
	opentracing.SetGlobalTracer(trace)

	neo4jClient := server.initNeo4J()

	loggingService := server.initLoggingService()

	messageService := server.initMessageService()

	connectionStore := server.initConnectionStore(neo4jClient, loggingService, messageService)

	connectionService := server.initConnectionService(connectionStore)

	connectionHandler := server.initConnectionHandler(connectionService)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(connectionService, replyPublisher, commandSubscriber)

	server.startGrpcServer(connectionHandler)
}

func (server *Server) initRegisterUserHandler(connectionService *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(connectionService, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initNeo4J() *neo4j.Driver {
	fmt.Println(fmt.Sprintf("%s://%s:%s", server.config.Neo4jUri, server.config.Neo4jHost, server.config.Neo4jPort))
	neo4jServer := fmt.Sprintf("%s://%s:%s", server.config.Neo4jUri, server.config.Neo4jHost, server.config.Neo4jPort)

	client, err := persistence.GetClient(neo4jServer, server.config.Neo4jUsername, server.config.Neo4jPassword)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initConnectionStore(client *neo4j.Driver, loggingService pbLogg.LoggingServiceClient, messageService pbMessage.MessageServiceClient) domain.ConnectionStore {
	store := persistence.NewConnectionDBStore(client, loggingService, messageService)
	store.Init()
	return store
}

func (server *Server) initConnectionService(store domain.ConnectionStore) *application.ConnectionService {
	return application.NewConnectionService(store, server.config)
}

func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewConnectionHandler(service)
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	creds, err := credentials.NewServerTLSFromFile("./certificates/connection_service.crt", "./certificates/connection_service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	connection.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initLoggingService() pbLogg.LoggingServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.LoggingHost, server.config.LoggingPort)
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Logging service: %v", err)
	}
	return pbLogg.NewLoggingServiceClient(conn)
}

func (server *Server) initMessageService() pbMessage.MessageServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.MessageHost, server.config.MessagePort)
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Message service: %v", err)
	}
	return pbMessage.NewMessageServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}
