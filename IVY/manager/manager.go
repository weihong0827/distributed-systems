package manager

import (
	"sync"

	pb "Ivy/pb"
)

type Page struct {
	content string
	owner   int64
	copySet []int64
}

type Manager struct {
	pb.UnimplementedManagerServiceServer
	pages              map[string]*Page
	requests           []*pb.WriteRequest
	waitingForResponse bool
	nodes              map[int64]string
	mu                 sync.Mutex
}

func NewManager(nodes map[int64]string) *Manager {
	return &Manager{
		pages: make(map[string]*Page),
		nodes: nodes,
	}
}
