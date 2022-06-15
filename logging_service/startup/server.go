package startup

import (
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/startup/config"
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

	loggingStore := server.initLoggingStore()

	loggingService := server.initLoggingService(loggingStore)

	loggingHandler := server.initLoggingHandler(loggingService)

	server.startGrpcServer(loggingHandler)
}

func (server *Server) initLoggingStore() persistence.LoggingStore {
	return persistence.NewLoggingDbStore(server.config)
}

func (server *Server) initLoggingService(store persistence.LoggingStore) *application.LoggingService {
	return application.NewLoggingService(store)
}

func (server *Server) initLoggingHandler(service *application.LoggingService) *api.LoggingHandler {
	return api.NewLoggingHandler(service)
}

func (server *Server) startGrpcServer(profileHandler *api.LoggingHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	loggingS.RegisterLoggingServiceServer(grpcServer, profileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
