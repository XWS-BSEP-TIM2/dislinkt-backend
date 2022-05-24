package startup

import (
	"context"
	"fmt"

	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/domain"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/auth_service/startup/config"
	auth "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/auth_service"
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
	credentialStore := server.initCredentialStore(mongoClient)
	passwordlessTokenStore := server.initPasswordlessTokenStore(mongoClient)
	emailService := server.initEmailService()
	authService := server.initAuthService(credentialStore, emailService)
	passwordlessLoginService := server.initPasswordlessLoginService(passwordlessTokenStore, emailService)

	authHandler := server.initAuthHandler(authService, passwordlessLoginService)

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

func (server *Server) initAuthService(store domain.UserStore, emailService *application.EmailService) *application.AuthService {
	profileServiceEndpoint := fmt.Sprintf("%s:%s", server.config.ProfileServiceHost, server.config.ProfileServicePort)
	return application.NewAuthService(store, profileServiceEndpoint, emailService)
}

func (server *Server) initAuthHandler(service *application.AuthService, passwordlessService *application.PasswordlessTokenService) *api.AuthHandler {
	return api.NewAuthHandler(service, passwordlessService)
}

func (server *Server) startGrpcServer(authHandler *api.AuthHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
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
