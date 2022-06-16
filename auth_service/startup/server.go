package startup

import (
	"context"
	"crypto/tls"
	"fmt"

	pbLogg "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"google.golang.org/grpc/credentials"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/startup/config"
	auth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
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
	credentialStore := server.initCredentialStore(mongoClient)
	passwordlessTokenStore := server.initPasswordlessTokenStore(mongoClient)
	apiTokenStore := server.initApiTokenStore(mongoClient)

	loggingService := server.initLoggingService()

	emailService := server.initEmailService()
	authService := server.initAuthService(credentialStore, emailService, loggingService)
	apiTokenService := server.initApiTokenService(apiTokenStore)
	passwordlessLoginService := server.initPasswordlessLoginService(passwordlessTokenStore, emailService)

	authHandler := server.initAuthHandler(authService, passwordlessLoginService, apiTokenService, loggingService)

	server.startGrpcServer(authHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.UserDBHost, server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initCredentialStore(client *mongo.Client) domain.UserStore {
	store := persistence.NewCredentialsMongoDBStore(client)
	store.DeleteAll(context.TODO())
	for _, user := range users {
		err, _ := store.Insert(context.TODO(), user)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initPasswordlessTokenStore(client *mongo.Client) persistence.PasswordlessTokenMongoDBStore {
	store := persistence.NewPasswordlessTokenMongoDBStore(client)
	return store
}

func (server *Server) initAuthService(store domain.UserStore, emailService *application.EmailService, loggingService pbLogg.LoggingServiceClient) *application.AuthService {
	profileServiceEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileServiceHost, server.config.ProfileServicePort)
	return application.NewAuthService(store, profileServiceEndpoint, emailService, loggingService)
}

func (server *Server) initAuthHandler(service *application.AuthService, passwordlessService *application.PasswordlessTokenService, tokenService *application.ApiTokenService, loggingService pbLogg.LoggingServiceClient) *api.AuthHandler {
	return api.NewAuthHandler(service, passwordlessService, tokenService, loggingService)
}

func (server *Server) startGrpcServer(authHandler *api.AuthHandler) {
	creds, err := credentials.NewServerTLSFromFile("./certificates/auth_service.crt", "./certificates/auth_service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	auth.RegisterAuthServiceServer(grpcServer, authHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initEmailService() *application.EmailService {
	return application.NewEmailService(server.config.Email, server.config.PasswordEmail, server.config.ApiGatwayHost, server.config.ApiGatwayPort)
}

func (server *Server) initPasswordlessLoginService(store persistence.PasswordlessTokenMongoDBStore, service *application.EmailService) *application.PasswordlessTokenService {
	return application.NewPasswordlessTokenService(&store, service)
}

func (server *Server) initApiTokenStore(client *mongo.Client) persistence.ApiTokenMongoDBStore {
	store := persistence.NewApiTokenMongoDBStore(client)
	store.DeleteAllTokens(context.TODO())
	for _, apiToken := range apiTokens {
		err, _ := store.Insert(context.TODO(), apiToken)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store

}

func (server *Server) initApiTokenService(store persistence.ApiTokenMongoDBStore) *application.ApiTokenService {
	return application.NewApiTokenService(&store)
}

func (server *Server) initLoggingService() pbLogg.LoggingServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.LoggingHost, server.config.LoggingPort)
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway faild to start", "Failed to start")
		log.Fatalf("Failed to start gRPC connection to Logging service: %v", err)
	}
	return pbLogg.NewLoggingServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	return grpc.Dial(address, grpc.WithTransportCredentials(credentials.NewTLS(config)))
}
