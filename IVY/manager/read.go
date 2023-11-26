package manager

import (
	"context"
	"errors"
	"log"

	pb "Ivy/pb"

	"Ivy/utils"
)

func (m *Manager) Read(ctx context.Context, req *pb.ReadRequest) (*pb.Empty, error) {
	log.Printf("Manager receive read request: %v", req)
	m.mu.Lock()
	if _, ok := m.pages[req.Page]; !ok {
		m.mu.Unlock()
		return nil, errors.New("Page not found")
	}

	currentPage := m.pages[req.Page]
	currentPage.copySet = append(currentPage.copySet, req.Source)
	m.mu.Unlock()

	client, conn, err := utils.CreateNodeServiceClient(m.nodes[currentPage.owner])
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	_, err = client.ReadForward(context.Background(),
		&pb.ReadForwardRequest{
			Page: req.Page,
			Node: m.nodes[req.Source],
		})
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (m *Manager) ReadConfirmation(ctx context.Context, req *pb.ReadConfirmationRequest) (*pb.Empty, error) {
	log.Printf("Read confirmation received at manager for request: %v Read Complete", req)
	return &pb.Empty{}, nil
}
