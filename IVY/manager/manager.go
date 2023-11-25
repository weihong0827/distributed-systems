package manager

import (
	"context"
	"fmt"
	"sync"

	pb "Ivy/pb"

	"google.golang.org/grpc"
)

type Page struct {
	content string
	owner   int
	copySet []int
}

type Manager struct {
	pb.UnimplementedManagerServiceServer
	pages              map[string]Page
	requests           []*pb.WriteRequest
	waitingForResponse bool
	nodes              map[int]string
	mu                 sync.Mutex
}

func NewManager(nodes map[int]string) *Manager {
	return &Manager{
		pages: make(map[string]Page),
		nodes: nodes,
	}
}

func CreateGRPCConnection(address string) (pb.NodeServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("Connection error", err)
	}
	client := pb.NewNodeServiceClient(conn)
	return client, conn, nil
}

func (m *Manager) sendInvalidateRequest(node string, failedNodes *[]string) bool {
	client, conn, err := CreateGRPCConnection(node)
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
	failedNodes := []string{}
	for _, node := range m.pages[m.requests[0].Page].copySet {
		m.sendInvalidateRequest(m.nodes[node], &failedNodes)
	}
	for len(failedNodes) > 0 {
		if m.sendInvalidateRequest(failedNodes[0], &failedNodes) {
			failedNodes = failedNodes[1:]
		}
	}
}

func (m *Manager) ServingWrites() {
	for {
		if !m.waitingForResponse && len(m.requests) > 0 {
			m.mu.Lock()
			m.waitingForResponse = true
			m.HandleWrite()
			m.mu.Unlock()

		}
	}
}

func (m *Manager) WriteRequest(ctx context.Context, req *pb.WriteRequest) (*pb.Empty, error) {
	m.mu.Lock()
	m.requests = append(m.requests, req)
	m.mu.Unlock()

	return nil, nil
}
