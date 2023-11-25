package node

import (
	"context"
	"log"

	pb "Ivy/pb"
)

type Access int

const (
	READ Access = iota
	WRITE
)

type LocalPage struct {
	content string
	access  Access
}

type Node struct {
	pb.UnimplementedNodeServiceServer
	pages map[string]LocalPage
	id    int
	CM    string
}

func NewNode(id int, CM string) *Node {
	return &Node{
		pages: make(map[string]LocalPage),
		id:    id,
		CM:    CM,
	}
}

func (n *Node) Invalidate(context context.Context, request *pb.InvalidateRequest) (*pb.Empty, error) {
	log.Printf("node %d: invalidate page %s", n.id, request.Page)
	// delete(n.pages, request.Page)
	return &pb.Empty{}, nil
}
