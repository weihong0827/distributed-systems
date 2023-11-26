package node

import (
	"sync"

	pb "Ivy/pb"
)

type Access int

const (
	READ Access = iota
	WRITE
	NIL
)

type LocalPage struct {
	content string
	access  Access
}

type Node struct {
	pb.UnimplementedNodeServiceServer
	pages         map[string]*LocalPage
	id            int64
	CM            string
	mu            sync.Mutex
	pendingWrites map[string]*pb.InitWriteRequest
	waitChan      chan string
}

func NewNode(id int64, CM string) *Node {
	return &Node{
		pages:         make(map[string]*LocalPage),
		id:            id,
		CM:            CM,
		pendingWrites: make(map[string]*pb.InitWriteRequest),
		waitChan:      make(chan string),
	}
}
