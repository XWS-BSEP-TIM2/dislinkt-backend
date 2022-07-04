package startup

import (
	"context"
	"crypto/tls"
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	profile "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/profile_service"
	saga "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/saga/messaging/nats"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/profile_service/startup/config"
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

const (
	QueueGroup = "profile_service"
)

func (server *Server) Start() {
	trace, _ := tracer.Init("profile_service")
	opentracing.SetGlobalTracer(trace)

	mongoClient := server.initMongoClient()

	loggingService := server.initLoggingService()

	profileStore := server.initProfileStore(mongoClient, loggingService)

	profileService := server.initProfileService(profileStore)

	profileHandler := server.initProfileHandler(profileService)

	commandSubscriber := server.initSubscriber(server.config.RegisterUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.RegisterUserReplySubject)
	server.initRegisterUserHandler(profileService, replyPublisher, commandSubscriber)

	fmt.Println("Profile service started.")
	server.startGrpcServer(profileHandler)
}

func (server *Server) initRegisterUserHandler(authService *application.ProfileService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewRegisterUserCommandHandler(authService, publisher, subscriber)
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

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.ProfileDBHost, server.config.ProfileDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initProfileStore(client *mongo.Client, loggingService loggingS.LoggingServiceClient) persistence.ProfileStore {
	store := persistence.NewProfileMongoDbStore(client, loggingService)

	store.DeleteAll(context.TODO())
	for _, user := range users {
		err := store.Insert(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
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

func (server *Server) initProfileService(store persistence.ProfileStore) *application.ProfileService {
	return application.NewProfileService(store)
}

func (server *Server) initProfileHandler(service *application.ProfileService) *api.ProfileHandler {
	return api.NewProfileHandler(service)
}

func (server *Server) startGrpcServer(profileHandler *api.ProfileHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)
	profile.RegisterProfileServiceServer(grpcServer, profileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("./certificates/profile_service.crt", "./certificates/profile_service.key")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
