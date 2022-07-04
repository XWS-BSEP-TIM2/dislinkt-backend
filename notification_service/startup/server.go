package startup

import (
	"context"
	"crypto/tls"
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	notificationS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/notification_service/startup/config"
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
	trace, _ := tracer.Init("notification_service")
	opentracing.SetGlobalTracer(trace)

	mongoClient := server.initMongoClient()

	loggingService := server.initLoggingService()

	notificationStore := server.initNotificationStore(mongoClient)

	notificationService := server.initNotificationService(notificationStore, loggingService)

	notificationHandler := server.initNotificationHandler(notificationService)

	fmt.Println("Notification service started.")
	server.startGrpcServer(notificationHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.NotificationDBHost, server.config.NotificationDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initNotificationStore(client *mongo.Client) persistence.NotificationStore {
	store := persistence.NewNotificationMongoDbStore(client)

	store.DeleteAllNotifications(context.TODO())
	store.DeleteAllSettings(context.TODO())
	for _, notification := range notifications {
		err := store.Insert(context.TODO(), notification)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, setting := range userSettings {
		err := store.InsertSetting(context.TODO(), setting)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
}

func (server *Server) initNotificationService(store persistence.NotificationStore, loggingService loggingS.LoggingServiceClient) *application.NotificationService {
	return application.NewNotificationService(store, server.config, loggingService)
}

func (server *Server) initNotificationHandler(service *application.NotificationService) *api.NotificationHandler {
	return api.NewNotificationHandler(service)
}

func (server *Server) startGrpcServer(notificationHandler *api.NotificationHandler) {
	creds, err := credentials.NewServerTLSFromFile("./certificates/notification_service.crt", "./certificates/notification_service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	notificationS.RegisterNotificationServiceServer(grpcServer, notificationHandler)
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
