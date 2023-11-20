package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Request struct {
	timestamp int
	pid       int
}

type LamportQueue struct {
	mu           sync.Mutex
	queue        []Request
	clock        int
	inCS         bool
	processPid   int
	channels     map[int]chan Request
	releaseChans map[int]chan bool
	replyChans   map[int]chan bool
	ownRequest   *Request
}

func NewLamportQueue(
	pid int,
	channels map[int]chan Request,
	releaseChans, replyChans map[int]chan bool,
) *LamportQueue {
	return &LamportQueue{
		queue:        make([]Request, 0),
		clock:        0,
		processPid:   pid,
		inCS:         false,
		channels:     channels,
		releaseChans: releaseChans,
		replyChans:   replyChans,
	}
}

func (l *LamportQueue) RequestCS() {
	l.mu.Lock()
	l.clock++
	request := Request{
		timestamp: l.clock,
		pid:       l.processPid,
	}
	l.queue = append(l.queue, request)
	l.ownRequest = &request
	l.mu.Unlock()

	// Broadcast request to all processes
	l.mu.Lock()
	fmt.Printf("Process %d broadcasting its request at clock %d.\n", l.processPid, l.clock)
	l.mu.Unlock()
	for pid, ch := range l.channels {
		if pid != l.processPid {
			ch <- request
		}
	}

	repliesReceived := 0
	for repliesReceived < len(l.channels)-1 {
		<-l.replyChans[l.processPid]
		repliesReceived++
		fmt.Printf("Process %d received reply. Total replies: %d\n", l.processPid, repliesReceived)
	}

	// Ensure this process request is at the head of the queue before entering CS
	for {
		l.mu.Lock()
		if l.queue[0].pid == l.processPid {
			l.mu.Unlock()
			break
		}
		l.mu.Unlock()
		time.Sleep(10 * time.Millisecond) // Polling interval
	}
	fmt.Printf("Node %d Entering Critical Section\n", l.processPid)
}

func (l *LamportQueue) ReleaseCS() {
	l.mu.Lock()
	l.queue = l.queue[1:]
	l.inCS = false
	l.ownRequest = nil
	l.mu.Unlock()

	// Broadcast release to all processes
	fmt.Printf("Process %d releasing its request.\n", l.processPid)
	for pid, ch := range l.releaseChans {
		if pid != l.processPid {
			ch <- true
		}
	}
}
func Less(req1 Request, req2 Request) bool {
	if req1.timestamp == req2.timestamp {
		return req1.pid < req2.pid
	}
	return req1.timestamp < req2.timestamp
}

func (l *LamportQueue) ReceiveRequest(request Request) {
	l.mu.Lock()
	fmt.Printf(
		"Process %d received request from Process %d at local clock %d. \n",
		l.processPid,
		request.pid,
		l.clock,
	)
	l.clock = max(l.clock, request.timestamp) + 1
	l.queue = append(l.queue, request)
	l.sortQueue()

	if l.inCS || (l.ownRequest != nil && Less(*l.ownRequest, request)) {
		l.mu.Unlock()
		return
	}

	l.mu.Unlock()
	fmt.Printf("Process %d sending reply to Process %d\n", l.processPid, request.pid)

	// Send reply back to the requesting process
	l.replyChans[request.pid] <- true
}

func (l *LamportQueue) sortQueue() {
	sort.Slice(l.queue, func(i, j int) bool {
		if l.queue[i].timestamp == l.queue[j].timestamp {
			return l.queue[i].pid < l.queue[j].pid
		}
		return l.queue[i].timestamp < l.queue[j].timestamp
	})
}

func (l *LamportQueue) Listen() {
	for {
		select {
		case request := <-l.channels[l.processPid]:
			l.ReceiveRequest(request)
			// race condition here
		case <-l.releaseChans[l.processPid]:
			l.mu.Lock()
			l.queue = l.queue[1:]
			l.replyChans[l.processPid] <- true
			l.mu.Unlock()
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	nClient := 10
	channels := make(map[int]chan Request, nClient)
	releaseChans := make(map[int]chan bool, nClient)
	replyChans := make(map[int]chan bool, nClient)

	for i := 1; i <= nClient; i++ {
		channels[i] = make(chan Request, nClient)
		releaseChans[i] = make(chan bool, nClient)
		replyChans[i] = make(chan bool, nClient)
	}

	processes := make([]*LamportQueue, nClient)
	for i := 1; i <= nClient; i++ {
		processes[i-1] = NewLamportQueue(i, channels, releaseChans, replyChans)
		go processes[i-1].Listen()
	}

	for i := 1; i <= nClient; i++ {
		go func(pid int) {
			processes[pid-1].RequestCS()
			fmt.Printf("Process %d in critical section\n", pid)
			time.Sleep(3 * time.Second)
			processes[pid-1].ReleaseCS()
		}(i)
	}

	time.Sleep(time.Duration(nClient+10) * time.Second)
}
