package startup

import (
	"crypto/tls"
	"fmt"
	loggingS "github.com/XWS-BSEP-TIM2/dislinkt-backend/common/proto/logging_service"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/common/tracer"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/application"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM2/dislinkt-backend/logging_service/startup/config"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	myLogger "gopkg.in/natefinch/lumberjack.v2"
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
	trace, _ := tracer.Init("logging_service")
	opentracing.SetGlobalTracer(trace)

	loggingStore := server.initLoggingStore()

	loggingService := server.initLoggingService(loggingStore)

	loggingHandler := server.initLoggingHandler(loggingService)
	//loggingService.LoggSuccess(context.TODO(), &loggingS.LogRequest{ServiceName: "LOG_SERVICE", ServiceFunctionName: "Start", Description: "Testiramo da li ovo radi", IpAddress: "localHost", UserID: "rasti"})

	fmt.Println("Logging service started.")
	server.startGrpcServer(loggingHandler)
}

func (server *Server) initLoggingStore() persistence.LoggingStore {
	logg := &myLogger.Logger{
		Filename:   server.config.FilePath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     28,    // days
		Compress:   false, // disabled by default
	}
	return persistence.NewLoggingDbStore(server.config, logg)
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

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)
	loggingS.RegisterLoggingServiceServer(grpcServer, profileHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("./certificates/logging_service.crt", "./certificates/logging_service.key")
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
