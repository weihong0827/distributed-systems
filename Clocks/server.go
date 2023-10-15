package main

import (
	"fmt"
	"math/rand"
	"sync"

	. "github.com/weihong0827/clockUsage/types"
)

type Server struct {
	Clients         []*Client
	broadcast       chan *Message
	clock           int
	vectorClock     []int
	useVectoreClock bool
	clockMutex      sync.Mutex
	stopCh          chan bool
}

func NewServer(size int, useVectoreClock bool, stopCh chan bool) *Server {
	return &Server{
		Clients:         make([]*Client, 0),
		broadcast:       make(chan *Message),
		clock:           0,
		vectorClock:     make([]int, size),
		useVectoreClock: useVectoreClock,
		stopCh:          stopCh,
	}
}

func (s *Server) register(c *Client) {
	s.Clients = append(s.Clients, c)
}

func (s *Server) IncrementAndUpdateMsg(msg *Message) *Message {
	// increment clock for sending and dropping message

	s.clock++
	s.IncrementVectorClock(-1)

	// update the message clock
	newMsg := &Message{
		Msg:         msg.Msg,
		From:        msg.From,
		Clock:       s.clock,
		VectorClock: make([]int, len(s.vectorClock)), // Ensure correct length for the new vector clock.
	}

	// Copy the server's vector clock to the new message's vector clock.
	copy(newMsg.VectorClock, s.vectorClock)

	return newMsg
}

func (s *Server) Start() {
	for {
		select {
		case <-s.stopCh:
			close(s.broadcast)
		case msg, ok := <-s.broadcast:
			if !ok {
				return
			}
			// Receive message from client
			s.clockMutex.Lock()
			UpdateClockWithMsg(s, *msg, -1)
			s.clockMutex.Unlock()

			if s.useVectoreClock && CheckCausalViolation(s, *msg) {
				fmt.Printf("Server found causal violation with msg clock %v and server clock %v \n", msg.GetVectorClock, s.GetVectorClock)
			} else if s.useVectoreClock {
				fmt.Printf("Server received: %s ,at server vector clock %v\n", msg.Msg, s.vectorClock)
			} else {
				fmt.Printf("Server received: %s ,at server clock %d\n", msg.Msg, s.clock)
			}

			if rand.Intn(2) == 0 {
				s.clockMutex.Lock()
				s.IncrementAndUpdateMsg(msg)
				if s.useVectoreClock {
					fmt.Printf("Dropped:%s at server vector clock %v \n", msg.Msg, s.vectorClock)
				} else {
					fmt.Printf("Dropped:%s at server clock %d \n", msg.Msg, s.clock)
				}
				s.clockMutex.Unlock()
				// drop the message
				continue
			}
			for _, client := range s.Clients {
				if client.id == msg.From {
					continue
				}
				s.clockMutex.Lock()
				updatedMsg := s.IncrementAndUpdateMsg(msg)
				if s.useVectoreClock {
					fmt.Printf("Server sent:%s to client %d at server vector clock %v \n", updatedMsg.Msg, client.id, s.vectorClock)
				} else {
					fmt.Printf("Server sent:%s to client %d at server clock %d \n", updatedMsg.Msg, client.id, s.clock)
				}
				client.inbox <- *updatedMsg
				s.clockMutex.Unlock()
			}
		}
	}
}

func (s *Server) GetClock() int {
	return s.clock
}

func (s *Server) GetVectorClock() []int {
	return s.vectorClock
}

func (s *Server) SetVectorClock(vector []int) {
	s.vectorClock = vector
}

func (s *Server) IncrementVectorClock(index int) {
	if index == -1 {
		index = len(s.vectorClock) - 1
	}
	s.vectorClock[index] += 1
}

func (s *Server) SetClock(clock int) {
	s.clock = clock
}
