package startup

import (
	"context"
	"fmt"
	joboffer "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/job_offer_service/startup/config"
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
	jobOfferStore := server.initJobOfferStore(mongoClient)

	jobOfferService := server.initJobOfferService(jobOfferStore)

	jobOfferHandler := server.initJobOfferHandler(jobOfferService)

	fmt.Println("Job offer service started.")
	server.startGrpcServer(jobOfferHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.JobOfferDBHost, server.config.JobOfferDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initJobOfferStore(client *mongo.Client) persistence.JobOfferStore {
	store := persistence.NewJobOfferMongoDbStore(client)

	store.DeleteAll(context.TODO())
	for _, jobOffer := range jobOffers {
		err := store.Insert(context.TODO(), jobOffer)
		if err != nil {
			log.Fatal(err)
		}
	}

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
