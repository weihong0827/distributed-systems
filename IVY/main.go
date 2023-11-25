package main

import (
	"fmt"
	"net"

	"Ivy/manager"
	"Ivy/node"

	pb "Ivy/pb"

	"google.golang.org/grpc"
)

func main() {
	nClient := 5
	nodes := make(map[int]string, nClient)

	CMAddress := "localhost:5000"

	lis, err := net.Listen("tcp", CMAddress)
	if err != nil {
		fmt.Printf("failed to listen on %s: %v\n", CMAddress, err)
		return
	}

	// Create a server instance for each node and pass the central manager
	for i := 1; i <= nClient; i++ {
		nodes[i] = fmt.Sprintf("localhost:500%d", i)
		go serveNode(i, nodes[i], CMAddress) // Start each node in a separate goroutine
	}

	CM := manager.NewManager(nodes)
	go CM.ServingWrites()

	grpcServer := grpc.NewServer()
	pb.RegisterManagerServiceServer(grpcServer, CM)
	fmt.Printf("Central Manager serving on %s\n", CMAddress)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve on %s: %v\n", CMAddress, err)
	}
}

func serveNode(nodeID int, address string, CM string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("failed to listen on %s: %v\n", address, err)
		return
	}

	grpcServer := grpc.NewServer()

	// Assuming you have a way to register the node with the CM
	node := node.NewNode(nodeID, CM)
	pb.RegisterNodeServiceServer(grpcServer, node)

	fmt.Printf("Node %d serving on %s\n", nodeID, address)
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve on %s: %v\n", address, err)
	}
}
