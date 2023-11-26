package manager

import (
	"context"
	"fmt"
	"log"

	pb "Ivy/pb"
	"Ivy/utils"
)

func (m *Manager) sendInvalidateRequest(node string, failedNodes *[]string) bool {
	client, conn, err := utils.CreateNodeServiceClient(node)
	if err != nil {
		*failedNodes = append(*failedNodes, node)
		fmt.Println("Handle write error:" + err.Error())
		return false
	}
	_, err = client.Invalidate(
		context.Background(),
		&pb.InvalidateRequest{
			Page: m.requests[0].Page,
		})
	if err != nil {
		fmt.Println("Handle write invalidate error:" + err.Error())
	}
	conn.Close()
	return true
}

func (m *Manager) HandleWrite() {
	m.mu.Lock()
	currentRequest := m.requests[0]
	value, ok := m.pages[currentRequest.Page]
	if !ok {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	failedNodes := []string{}
	for _, node := range value.copySet {
		m.sendInvalidateRequest(m.nodes[node], &failedNodes)
	}
	for len(failedNodes) > 0 {
		if m.sendInvalidateRequest(failedNodes[0], &failedNodes) {
			failedNodes = failedNodes[1:]
		}
	}
	// Forward write to previous owner
	log.Printf("Forwarding write to %d", value.owner)
	client, conn, err := utils.CreateNodeServiceClient(m.nodes[value.owner])
	if err != nil {
		fmt.Println("Handle write error:" + err.Error())
		return
	}
	_, err = client.WriteForward(context.Background(), &pb.ForwardRequest{
		Page: currentRequest.Page,
		Node: m.nodes[currentRequest.Source],
	})
	if err != nil {
		fmt.Println("Handle write forward error:" + err.Error())
	}
	conn.Close()
	log.Print("exiting handle write")
}

func (m *Manager) ServingWrites() {
	for {
		m.mu.Lock()
		if !m.waitingForResponse && len(m.requests) > 0 {
			log.Printf("manager: serving write request for page %s", m.requests[0].Page)
			m.waitingForResponse = true
			m.mu.Unlock()
			m.HandleWrite()

		} else {
			m.mu.Unlock()
		}
	}
}

func (m *Manager) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	log.Printf("manager: write request for page %s", req.Page)
	m.mu.Lock()
	log.Print("append to request")
	m.requests = append(m.requests, req)
	log.Printf("requests: %v", m.requests)
	m.mu.Unlock()
	if _, ok := m.pages[req.Page]; !ok {
		return &pb.WriteResponse{
			ToWrite: true,
		}, nil
	}

	return &pb.WriteResponse{
		ToWrite: false,
	}, nil
}

func (m *Manager) WriteConfirmation(ctx context.Context, req *pb.WriteConfirmationRequest) (*pb.Empty, error) {
	log.Printf("manager: write confirmation for page %s", req.Page)

	m.mu.Lock()
	defer m.mu.Unlock()

	if page, ok := m.pages[req.Page]; ok {
		if page.owner == req.Source {
			log.Printf("page owner %d writing to page %s ", page.owner, req.Page)
			page.content = req.Content
			return &pb.Empty{}, nil
		}
	}

	for m.waitingForResponse && len(m.requests) > 0 && m.requests[0].Page == req.Page && m.requests[0].Source == req.Source {
		m.waitingForResponse = false
		m.requests = m.requests[1:]
	}
	m.pages[req.Page] = &Page{
		content: req.Content,
		owner:   req.Source,
		copySet: []int64{},
	}
	log.Printf("Current Page on CM: %v", m.pages[req.Page])
	return &pb.Empty{}, nil
}
