package manager

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	pb "Ivy/pb"
	"Ivy/utils"
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
	isReplica          bool
	currtHeadCM        int64
	isAlive            bool
	managerList        map[int64]string
	id                 int64
}

func NewManager(nodes map[int64]string, isReplica bool, isAlive bool, currtHeadCM int64, managerList map[int64]string, id int64) *Manager {
	return &Manager{
		pages:       make(map[string]*Page),
		nodes:       nodes,
		isReplica:   isReplica,
		isAlive:     isAlive,
		currtHeadCM: currtHeadCM,
		managerList: managerList,
		id:          id,
	}
}

func (m *Manager) Watch() {
	for {
		if !m.isAlive {
			continue
		}
		if m.isReplica {
			client, _, err := utils.CreateManagerServiceClient(m.managerList[m.currtHeadCM])
			if err != nil {
				log.Print(err)
				m.switchToBackup()
				continue
			}

			_, err = client.HealthCheck(context.Background(), &pb.Empty{})
			if err != nil {
				log.Print(err)
				m.InitiateLeader()
			} else {
				m.synchronizeState(client)
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func (m *Manager) Start(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	m.mu.Lock()
	m.isAlive = true
	m.mu.Unlock()
	if !m.isReplica {
		for _, value := range m.managerList {
			client, _, err := utils.CreateManagerServiceClient(value)
			if err != nil {
				log.Print(err)
				continue
			}
			m.synchronizeState(client)
		}
	}

	return &pb.Empty{}, nil
}

func (m *Manager) Stop(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	m.mu.Lock()
	m.isAlive = false
	m.mu.Unlock()
	return &pb.Empty{}, nil
}

func (m *Manager) HealthCheck(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	if !m.isAlive {
		return nil, errors.New("Manager is not alive")
	}
	return &pb.Empty{}, nil
}

func (m *Manager) InitiateLeader() {
	for key, value := range m.managerList {
		if key > m.id {
			err := performHealthCheck(value)
			if err != nil {
				break
			}
		}
	}
	m.mu.Lock()
	m.currtHeadCM = m.id
	m.mu.Unlock()
	m.switchToBackup()
}

func performHealthCheck(address string) error {
	client, conn, err := utils.CreateManagerServiceClient(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.HealthCheck(context.Background(), &pb.Empty{})
	return err
}

func (m *Manager) switchToBackup() {
	// braod cast to all nodes
	for _, addr := range m.nodes {
		nodeClient, nodeConn, err := utils.CreateNodeServiceClient(addr)
		if err != nil {
			log.Print(err)
		}
		defer nodeConn.Close()
		_, err = nodeClient.SwitchCM(context.Background(), &pb.SwitchCMRequest{CMAddr: m.managerList[m.currtHeadCM]})
		if err != nil {
			log.Print(err)
		}
	}
	// braod cast to all managers
	for id, addr := range m.managerList {
		if id == m.id {
			continue
		}
		managerClient, managerConn, err := utils.CreateManagerServiceClient(addr)
		if err != nil {
			log.Print(err)
		}
		defer managerConn.Close()
		_, err = managerClient.SwitchCM(context.Background(), &pb.SwitchCMRequest{CMAddr: m.managerList[m.currtHeadCM]})
		if err != nil {
			log.Print(err)
		}
	}
}

func (m *Manager) SwitchCM(ctx context.Context, req *pb.SwitchCMRequest) (*pb.Empty, error) {
	if !m.isAlive {
		return nil, errors.New("Manager is not alive")
	}
	m.mu.Lock()
	m.currtHeadCM = m.id
	m.mu.Unlock()
	return &pb.Empty{}, nil
}

func (m *Manager) GetState(ctx context.Context, req *pb.Empty) (*pb.State, error) {
	if !m.isAlive || m.currtHeadCM != m.id {
		return nil, errors.New("get state failed")
	}
	pages := make(map[string]*pb.Page)
	for key, value := range m.pages {
		pages[key] = &pb.Page{
			Content: value.content,
			Owner:   value.owner,
			CopySet: value.copySet,
		}
	}
	return &pb.State{
		Pages:         pages,
		WriteRequests: m.requests,
	}, nil
}

func (m *Manager) synchronizeState(client pb.ManagerServiceClient) {
	state, err := client.GetState(context.Background(), &pb.Empty{})
	if err != nil {
		log.Print("Failed to synchronize state: ", err)
		return
	}
	log.Printf("Node %d synchronizing state get %v", m.id, state.Pages)

	m.mu.Lock()
	for key, value := range state.Pages {
		m.pages[key] = &Page{
			content: value.Content,
			owner:   value.Owner,
			copySet: value.CopySet,
		}
	}
	m.requests = state.WriteRequests
	m.mu.Unlock()
}
