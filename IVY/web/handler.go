package web

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	pb "Ivy/pb"
	"Ivy/utils"
)

func HandleWriteRequest(w http.ResponseWriter, r *http.Request, clients map[int64]string) {
	query := r.URL.Query()
	nodeID, _ := strconv.ParseInt(query.Get("node"), 10, 64)

	if clientAddr, ok := clients[nodeID]; ok {
		fmt.Fprintf(w, "Write request sent to node %d\n", nodeID)
		client, conn, err := utils.CreateNodeServiceClient(clientAddr)
		if err != nil {
			http.Error(w, "Error creating client", http.StatusInternalServerError)
		}

		_, err = client.InitWrite(context.Background(), &pb.InitWriteRequest{
			Page:    query.Get("page"),
			Content: query.Get("content"),
		})
		if err != nil {
			http.Error(w, "Error sending request", http.StatusInternalServerError)
		}

		conn.Close()

	} else {
		http.Error(w, "Invalid node ID", http.StatusBadRequest)
	}
}

func HandleReadRequest(w http.ResponseWriter, r *http.Request, clients map[int64]string) {
	query := r.URL.Query()
	nodeID, _ := strconv.ParseInt(query.Get("node"), 10, 64)

	if clientAddr, ok := clients[nodeID]; ok {
		fmt.Fprintf(w, "Read request sent to node %d\n", nodeID)
		client, conn, err := utils.CreateNodeServiceClient(clientAddr)
		if err != nil {
			http.Error(w, "Error creating client", http.StatusInternalServerError)
		}

		resp, err := client.InitRead(context.Background(), &pb.InitReadRequest{
			Page: query.Get("page"),
		})
		if err != nil {
			errMsg := fmt.Sprintf("Error sending request: %v", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		conn.Close()
		fmt.Fprintf(w, "The content for page %s is: %s", query.Get("page"), resp.Content)

	} else {
		http.Error(w, "Invalid node ID", http.StatusBadRequest)
	}
}
