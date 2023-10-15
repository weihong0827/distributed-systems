package main

import (
	"fmt"
	"sync"
	"time"

	. "github.com/weihong0827/clockUsage/types"
)

type ClientInterface interface {
	Send()
	Receive()
}
type Client struct {
	id             int
	inbox          chan Message
	toServer       chan *Message
	clock          int
	vectorClock    []int
	useVectorClock bool
	clockMutex     sync.Mutex
	StopCh         chan bool
	wg             sync.WaitGroup
	ReceivedMsgs   []Message
}

func NewClient(id int, toServer chan *Message, size int, useVectorClock bool) *Client {
	return &Client{
		id:             id,
		inbox:          make(chan Message, 10),
		toServer:       toServer,
		clock:          0,
		vectorClock:    make([]int, size),
		useVectorClock: useVectorClock,
	}
}

func (c *Client) Send() {
	defer c.wg.Done()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c.StopCh:
			return
		case <-ticker.C:
			// periodically send a message to the server
			var msg string
			c.clockMutex.Lock()
			c.clock++
			c.IncrementVectorClock(c.id)
			if c.useVectorClock {
				msg = fmt.Sprintf("Hello from client %d at time %d", c.id, c.GetVectorClock()[c.id])
			} else {
				msg = fmt.Sprintf("Hello from client %d at time %d", c.id, c.clock)
			}
			c.toServer <- &Message{
				Msg:         msg,
				From:        c.id,
				Clock:       c.clock,
				VectorClock: c.vectorClock,
			}
			c.clockMutex.Unlock()
		}
	}
}

func (c *Client) Receive() {
	defer c.wg.Done()
	for {
		select {
		case <-c.StopCh:
			return
			// receive messages from the server
		case msg := <-c.inbox:
			c.clockMutex.Lock()
			UpdateClockWithMsg(c, msg, c.id)
			c.clockMutex.Unlock()
			c.ReceivedMsgs = append(c.ReceivedMsgs, msg)
			if c.useVectorClock && CheckCausalViolation(c, msg) {
				fmt.Printf("Client %d found causal violation with msg clock %v and client clock %v \n", c.id, msg.GetVectorClock(), c.GetVectorClock())
			} else if c.useVectorClock {
				fmt.Printf("Client %d received:%s with msg clock %v and client clock %v \n", c.id, msg.Msg, msg.GetVectorClock(), c.GetVectorClock())
			} else {
				fmt.Printf("Client %d received: %s ,at cleint clock %d\n", c.id, msg.Msg, c.GetClock())
			}

		}
	}
}

func (c *Client) GetClock() int {
	return c.clock
}

func (c *Client) GetVectorClock() []int {
	return c.vectorClock
}

func (c *Client) SetVectorClock(vector []int) {
	c.vectorClock = vector
}

func (c *Client) IncrementVectorClock(index int) {
	if index == -1 {
		index = len(c.vectorClock) - 1
	}
	c.vectorClock[index] += 1
}

func (c *Client) SetClock(clock int) {
	c.clock = clock
}
