package services

import (
	//catalogue "github.com/tamararankovic/microservices_demo/common/proto/catalogue_service"
	//ordering "github.com/tamararankovic/microservices_demo/common/proto/ordering_service"
	//shipping "github.com/tamararankovic/microservices_demo/common/proto/shipping_service"
	post "github.com/tamararankovic/microservices_demo/common/proto/post_service"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewPostClient(address string) post.PostServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return post.NewPostServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
