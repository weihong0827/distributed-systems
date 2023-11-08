package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type MessageType int

const (
	Release MessageType = iota
	Rescind
	RescindReply
)

func (m MessageType) String() string {
	return [...]string{"Release", "Rescind", "RescindReply"}[m]
}

type Message struct {
	msg MessageType
	pid int
}

type Request struct {
	timestamp  int
	pid        int
	replyChann chan int
}
type Node struct {
	pid         int
	votedFor    *Request
	clock       int
	mu          sync.Mutex
	queue       []Request
	inCS        bool
	votedFrom   []int
	requestChan map[int]chan Request
	messageChan map[int]chan Message
}

func (n *Node) RequestCS() {
	n.mu.Lock()
	fmt.Printf("Node %d requesting CS at local clock %d\n", n.pid, n.clock)
	n.clock++
	ts := n.clock
	n.mu.Unlock()

	request := Request{
		timestamp:  ts,
		pid:        n.pid,
		replyChann: make(chan int),
	}

	for id, ch := range n.requestChan {
		if id == n.pid {
			continue
		}
		ch <- request
	}

	targetNum := len(n.messageChan)/2 + 1
	replies := 0
	for replies < targetNum {
		sourceId := <-request.replyChann
		fmt.Printf("Node %d receive reply from node node %d but not updated yet\n", n.pid, sourceId)
		n.mu.Lock()
		n.votedFrom = append(n.votedFrom, sourceId)
		replies = len(n.votedFrom)
		fmt.Printf(
			"Node %d receive reply from node %d total replies: %d\n",
			n.pid,
			sourceId,
			replies,
		)
		n.mu.Unlock()
	}

	n.mu.Lock()
	n.inCS = true
	n.mu.Unlock()
	fmt.Printf("Node %d Entering CS\n", n.pid)

}
func (n *Node) ReleaseCS() {
	n.mu.Lock()
	defer n.mu.Unlock()
	releaseMsg := Message{
		msg: Release,
		pid: n.pid,
	}
	for _, id := range n.votedFrom {
		fmt.Printf("Node %d releasing vote to %d\n", n.pid, id)
		n.messageChan[id] <- releaseMsg
	}
	n.votedFrom = nil
	n.inCS = false

	if n.votedFor == nil {
		if len(n.queue) == 0 {
			return
		}
		n.sortQueue()
		n.sendReply()

	}

	fmt.Printf("Node %d exited CS\n", n.pid)
}

func Less(req1 Request, req2 Request) bool {

	if req1.timestamp == req2.timestamp {
		return req1.pid < req2.pid
	}
	return req1.timestamp < req2.timestamp
}

func remove(slice []int, value int) []int {
	for i, v := range slice {
		if v == value {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (n *Node) handleRequest(request Request) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.clock = max(n.clock, request.timestamp) + 1
	n.queue = append(n.queue, request)
	if n.votedFor == nil {
		fmt.Printf("Node %d vote for node %d\n", n.pid, request.pid)
		request.replyChann <- n.pid
		for i := range n.queue {
			if n.queue[i].pid == request.pid && n.queue[i].timestamp == request.timestamp {
				n.votedFor = &n.queue[i]
				break
			}
		}
	} else {
		fmt.Printf("Node %d already voted, checking for priority..\n", n.pid)
		fmt.Printf("votedFor: %d %d, incoming %d %d\n", n.votedFor.pid, n.votedFor.timestamp, request.pid, request.timestamp)
		if Less(request, *n.votedFor) {
			rescindMessage := Message{
				msg: Rescind,
				pid: n.pid,
			}
			fmt.Printf("Node %d sending rescind message to node %d\n", n.pid, n.votedFor.pid)
			n.messageChan[n.votedFor.pid] <- rescindMessage
		}
	}
}
func (n *Node) sendReply() {
	nextRequest := n.queue[0]
	n.votedFor = &nextRequest
	fmt.Printf("Node %d vote for node %d\n", n.pid, nextRequest.pid)
	nextRequest.replyChann <- n.pid
	fmt.Println("Sent")

}

func (n *Node) handleMessage(message Message) {

	n.mu.Lock()
	defer n.mu.Unlock()

	n.clock++
	switch message.msg {

	case Rescind:
		if n.inCS {
			return
		}
		n.votedFrom = remove(n.votedFrom, message.pid)
		releaseMsg := Message{
			msg: RescindReply,
			pid: n.pid,
		}
		n.messageChan[message.pid] <- releaseMsg
	// receive release vote
	case Release:
		if len(n.queue) > 0 {
			// Check if the first request in the queue is the one we have voted for
			fmt.Printf("queue at node %d is %v\n", n.pid, n.queue)
			n.queue = n.queue[1:] // Remove the first request
			// If there's a next request, vote for it
			if len(n.queue) > 0 {
				n.sortQueue()
				n.sendReply()
			} else {
				// No requests left to vote for
				n.votedFor = nil
			}
		} else {
			// Queue is empty, no need to vote for anything
			n.votedFor = nil
		}
	case RescindReply:
		// sort queue and reply
		n.votedFor = nil
		if n.inCS {
			return
		}

		if len(n.queue) > 0 {
			n.sortQueue()
			n.sendReply()
		}
	}

}
func (n *Node) Listen() {
	for {
		select {
		case request := <-n.requestChan[n.pid]:
			fmt.Printf("Receive CS request from %d at node %d\n", request.pid, n.pid)
			n.handleRequest(request)
		case message := <-n.messageChan[n.pid]:
			fmt.Printf("Receive message %s from %d at node %d\n", message.msg, message.pid, n.pid)
			n.handleMessage(message)

		}
	}

}

func NewNode(pid int, messageChan map[int]chan Message, requestChan map[int]chan Request) *Node {
	return &Node{
		pid:         pid,
		messageChan: messageChan,
		requestChan: requestChan,
	}
}

func (n *Node) sortQueue() {
	sort.Slice(n.queue, func(i, j int) bool {
		return Less(n.queue[i], n.queue[j])
	})
}

func main() {
	nClient := 5
	messageChan := make(map[int]chan Message, nClient)
	requestChan := make(map[int]chan Request, nClient)
	for i := 0; i < nClient; i++ {
		messageChan[i] = make(chan Message, nClient)
		requestChan[i] = make(chan Request, nClient)
	}

	nodes := make([]*Node, nClient)
	for i := 0; i < nClient; i++ {
		nodes[i] = NewNode(i, messageChan, requestChan)
		go nodes[i].Listen()
	}

	for i := 0; i < nClient; i++ {
		go func(pid int) {
			nodes[pid].RequestCS()
			fmt.Println("-------------------------------------")
			fmt.Printf("Process %d in critical section\n", pid)
			fmt.Println("-------------------------------------")
			time.Sleep(3 * time.Second)
			nodes[pid].ReleaseCS()
		}(i)
	}
	time.Sleep(35 * time.Second)
}
