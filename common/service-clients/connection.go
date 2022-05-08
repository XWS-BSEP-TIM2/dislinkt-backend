package service_clients

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
