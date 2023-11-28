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
	var nClient int64 = 5
	var nReplica int64 = 3
	nodes := make(map[int64]string, nClient)
	managers := make(map[int64]string, nReplica)

	primaryCM := "localhost:4000"

	// Create a server instance for each node and pass the central manager
	var i int64
	for i = 1; i <= nClient; i++ {
		nodes[i] = fmt.Sprintf("localhost:500%d", i)
		go serveNode(i, nodes[i], primaryCM) // Start each node in a separate goroutine
	}
	for i = 0; i < nReplica; i++ {
		managers[i] = fmt.Sprintf("localhost:400%d", i)
	}
	for id, addr := range managers {
		if id == 0 {
			go serveCM(primaryCM, addr, manager.NewManager(nodes, false, true, 0, managers, id))
		} else {
			go serveCM(primaryCM, addr, manager.NewManager(nodes, true, true, 0, managers, id))
		}
	}

	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWriteRequest(w, r, nodes)
	})
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleReadRequest(w, r, nodes)
	})
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleManagerState(w, r, managers, handler.STOP)
	})
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleManagerState(w, r, managers, handler.START)
	})

	http.ListenAndServe(":8080", nil)
}

func serveCM(primaryCM string, CMAddress string, CM *manager.Manager) {
	lis, err := net.Listen("tcp", CMAddress)
	if err != nil {
		fmt.Printf("failed to listen on %s: %v\n", CMAddress, err)
		return
	}
	go CM.ServingWrites()
	go CM.Watch(primaryCM)

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
