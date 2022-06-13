package startup

import (
	"context"
	"fmt"
	messageS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/message_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/message_service/startup/config"
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

	profileStore := server.initMessageStore(mongoClient)

	profileService := server.initMessageService(profileStore)

	profileHandler := server.initMessageHandler(profileService)

	server.startGrpcServer(profileHandler)
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
		err := store.Insert(context.TODO(), chat)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
}

func (server *Server) initMessageService(store persistence.MessageStore) *application.MessageService {
	return application.NewMessageService(store, server.config)
}

func (server *Server) initMessageHandler(service *application.MessageService) *api.MessageHandler {
	return api.NewMessageHandler(service)
}

func (server *Server) startGrpcServer(profileHandler *api.MessageHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	messageS.RegisterMessageServiceServer(grpcServer, profileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
