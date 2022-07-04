package startup

import (
	"crypto/tls"
	"fmt"
	joboffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	pbLogg "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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
	neo4jClient := server.initNeo4J()

	loggingService := server.initLoggingService()

	jobOfferStore := server.initJobOfferStore(neo4jClient, loggingService)

	jobOfferService := server.initJobOfferService(jobOfferStore)

	jobOfferHandler := server.initJobOfferHandler(jobOfferService)

	fmt.Println("Job offer service started.")
	server.startGrpcServer(jobOfferHandler)
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

func (server *Server) initLoggingService() pbLogg.LoggingServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.LoggingHost, server.config.LoggingPort)
	conn, err := getConnection(address)
	if err != nil {
		fmt.Println("Gateway failed to start", "Failed to start")
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

func (server *Server) initJobOfferStore(driver *neo4j.Driver, loggingService pbLogg.LoggingServiceClient) persistence.JobOfferStore {
	store := persistence.NewJobOfferDbStore(driver, loggingService)
	store.Init()
	return store
}

func (server *Server) initJobOfferService(store persistence.JobOfferStore) *application.JobOfferService {
	return application.NewJobOfferService(store)
}

func (server *Server) initJobOfferHandler(service *application.JobOfferService) *api.JobOfferHandler {
	return api.NewJobOfferHandler(service)
}

func (server *Server) startGrpcServer(profileHandler *api.JobOfferHandler) {
	creds, err := credentials.NewServerTLSFromFile("./certificates/job_offer_service.crt", "./certificates/job_offer_service.key")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	joboffer.RegisterJobOfferServiceServer(grpcServer, profileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
