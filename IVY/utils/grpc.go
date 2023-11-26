package utils

import (
	"fmt"

	"google.golang.org/grpc/credentials/insecure"

	pb "Ivy/pb"

	"google.golang.org/grpc"
)

func CreateGRPCConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Connection error", err)
	}
	return conn, nil
}

func CreateManagerServiceClient(address string) (pb.ManagerServiceClient, *grpc.ClientConn, error) {
	conn, err := CreateGRPCConnection(address)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewManagerServiceClient(conn), conn, nil
}

func CreateNodeServiceClient(address string) (pb.NodeServiceClient, *grpc.ClientConn, error) {
	conn, err := CreateGRPCConnection(address)
	if err != nil {
		return nil, nil, err
	}
	return pb.NewNodeServiceClient(conn), conn, nil
}
