package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"Ivy/manager"
	"Ivy/node"
	handler "Ivy/web"

	pb "Ivy/pb"

	"google.golang.org/grpc"
)

func main() {
	var nClient int64 = 10
	nodes := make(map[int64]string, nClient)

	CMAddress := "localhost:4000"

	// Create a server instance for each node and pass the central manager
	var i int64
	for i = 1; i <= nClient; i++ {
		nodes[i] = fmt.Sprintf("localhost:500%d", i)
		go serveNode(i, nodes[i], CMAddress) // Start each node in a separate goroutine
	}
	go serveCM(CMAddress, nodes)

	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWriteRequest(w, r, nodes)
	})
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleReadRequest(w, r, nodes)
	})
	http.ListenAndServe(":8080", nil)
}

func serveCM(CMAddress string, nodes map[int64]string) {
	lis, err := net.Listen("tcp", CMAddress)
	if err != nil {
		fmt.Printf("failed to listen on %s: %v\n", CMAddress, err)
		return
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

func serveNode(nodeID int64, address string, CM string) {
	time.Sleep(2 * time.Second)
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
