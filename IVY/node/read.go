package node

import (
	"context"
	"errors"
	"log"
	"time"

	pb "Ivy/pb"
	"Ivy/utils"
)

func (n *Node) InitRead(ctx context.Context, req *pb.InitReadRequest) (*pb.InitReadResponse, error) {
	log.Print("=========================================================================")
	log.Printf("Node %d InitRead request: %v", n.id, req)

	// If we have the page, return it
	if val, ok := n.pages[req.Page]; ok && val.access != NIL {
		log.Printf("Node %d has page %s", n.id, req.Page)
		return &pb.InitReadResponse{Content: n.pages[req.Page].content}, nil
	}

	// Otherwise, send a request to the manager
	client, conn, err := utils.CreateManagerServiceClient(n.CM)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	go client.Read(context.Background(), &pb.ReadRequest{Page: req.Page, Source: n.id})
	for {
		select {
		case page := <-n.readWaitChan:
			log.Printf("Node %d received page: %v", n.id, page)
			n.mu.Lock()
			n.pages[page.page] = &LocalPage{content: page.content, access: READ}
			n.mu.Unlock()

			return &pb.InitReadResponse{Content: page.content}, nil
		case <-time.After(2 * time.Second):
			return nil, errors.New("Timeout")
		}
	}
}

func (n *Node) ReadForward(ctx context.Context, req *pb.ReadForwardRequest) (*pb.Empty, error) {
	log.Printf("Node %d forwarding read request for page: %s to Node %s", n.id, req.Page, req.Node)
	n.mu.Lock()
	if _, ok := n.pages[req.Page]; !ok {
		n.mu.Unlock()
		return nil, errors.New("page not found")
	}

	currentPage := n.pages[req.Page]
	n.mu.Unlock()

	client, conn, err := utils.CreateNodeServiceClient(req.Node)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	_, err = client.SendContent(ctx, &pb.SendContentRequest{
		Page:    req.Page,
		Content: currentPage.content,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (n *Node) SendContent(ctx context.Context, req *pb.SendContentRequest) (*pb.Empty, error) {
	log.Printf("Node %d received content of request: %v", n.id, req)

	n.readWaitChan <- ReadPage{content: req.Content, page: req.Page}

	client, conn, err := utils.CreateManagerServiceClient(n.CM)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, err = client.ReadConfirmation(context.Background(), &pb.ReadConfirmationRequest{
		Page:   req.Page,
		Source: n.id,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
