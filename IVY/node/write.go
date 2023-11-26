package node

import (
	"context"
	"errors"
	"log"
	"time"

	pb "Ivy/pb"
	"Ivy/utils"
)

func (n *Node) Invalidate(context context.Context, request *pb.InvalidateRequest) (*pb.Empty, error) {
	log.Printf("node %d: invalidate page %s", n.id, request.Page)
	n.mu.Lock()
	if page, exists := n.pages[request.Page]; exists {
		page.access = NIL
	}
	n.mu.Unlock()
	return &pb.Empty{}, nil
}

func (n *Node) InitWrite(context context.Context, request *pb.InitWriteRequest) (*pb.Empty, error) {
	log.Printf("node %d: initiate write for page:%s with content %s", n.id, request.Page, request.Content)

	if page, ok := n.pages[request.Page]; ok {
		if page.access == WRITE {
			err := n.SendWriteConfirmation(context, &pb.WriteConfirmationRequest{
				Page:    request.Page,
				Content: request.Content,
				Source:  n.id,
			})
			if err != nil {
				return nil, err
			}
			n.UpdatePageContent(request.Page, request.Content, WRITE)
			return &pb.Empty{}, nil
		}
	}

	client, conn, err := utils.CreateManagerServiceClient(n.CM)
	if err != nil {
		return nil, err
	}

	resp, err := client.Write(context, &pb.WriteRequest{
		Page:   request.Page,
		Source: n.id,
	})
	if err != nil {
		log.Printf("node %d: error sending write request to manager: %s", n.id, err.Error())
		return nil, err
	}
	conn.Close()
	if resp.ToWrite {
		n.SendWriteConfirmation(context, &pb.WriteConfirmationRequest{
			Page:    request.Page,
			Content: request.Content,
			Source:  n.id,
		})
	} else {
		n.pendingWrites[request.Page] = request
		for {
			select {
			case page := <-n.waitChan:
				if page == request.Page {
					delete(n.pendingWrites, page)
					n.UpdatePageContent(request.Page, request.Content, WRITE)
					return &pb.Empty{}, nil
				}
			case <-time.After(10 * time.Second):
				return nil, errors.New("Timeout")
			}
		}
	}
	return &pb.Empty{}, nil
}

func (n *Node) UpdatePageContent(page string, content string, permission Access) {
	n.pages[page] = &LocalPage{
		content: content,
		access:  permission,
	}
	log.Printf("Current Page on node: %v", n.pages[page])
}

func (n *Node) SendWriteConfirmation(context context.Context, request *pb.WriteConfirmationRequest) error {
	log.Printf("node %d: send write confirmation for page %s", n.id, request.Page)
	client, conn, err := utils.CreateManagerServiceClient(n.CM)
	if err != nil {
		return err
	}
	_, err = client.WriteConfirmation(context, request)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (n *Node) WriteForward(context context.Context, request *pb.ForwardRequest) (*pb.Empty, error) {
	log.Printf("node %d: write forward for page %s", n.id, request.Page)
	client, conn, err := utils.CreateNodeServiceClient(request.Node)
	if err != nil {
		return nil, err
	}
	// Invalidate page on node
	n.Invalidate(context, &pb.InvalidateRequest{
		Page: request.Page,
	})

	_, err = client.Send(context, &pb.SendRequest{
		Page: request.Page,
	})
	if err != nil {
		return nil, err
	}
	conn.Close()
	return &pb.Empty{}, nil
}

func (n *Node) Send(context context.Context, request *pb.SendRequest) (*pb.Empty, error) {
	log.Printf("node %d: received request from previous owner for page %s", n.id, request.Page)
	client, conn, err := utils.CreateManagerServiceClient(n.CM)
	if err != nil {
		return nil, err
	}
	currPendingWrite := n.pendingWrites[request.Page]
	_, err = client.WriteConfirmation(context, &pb.WriteConfirmationRequest{
		Page:    request.Page,
		Content: currPendingWrite.Content,
		Source:  n.id,
	})
	if err != nil {
		return nil, err
	}
	conn.Close()
	log.Printf("node %d: sent write confirmation for page %s", n.id, request.Page)
	n.waitChan <- request.Page

	return &pb.Empty{}, nil
}
