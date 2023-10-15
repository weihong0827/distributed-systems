package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Node struct {
	id                     int
	data                   int // This is the data structure to be synchronized
	isCoordinator          bool
	coordinator            *Node
	isRunning              bool // Indicates if the node's routine is running
	nodes                  []*Node
	mu                     sync.Mutex
	electionMutex          sync.Mutex // prevent multiple simultaneous elections
	failDuringAnnouncement bool
	failDuringElection     bool
}

func (n *Node) run() {
	for {
		n.mu.Lock()
		if !n.isRunning {
			n.mu.Unlock()
			return
		}
		n.mu.Unlock()

		// Simulate periodic activity
		time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)

		if n.isCoordinator {
			continue // if it's the coordinator, no need to initiate an election
		}

		n.mu.Lock()
		coordinator := n.getCoordinator()
		n.mu.Unlock()

		success := n.syncData(coordinator)
		if !success {
			fmt.Printf("Node %d detected no coordinator\n", n.id)
			n.startElection()
		}
	}
}

func (n *Node) getCoordinator() *Node {
	return n.coordinator
}

func (n *Node) startElection() {
	n.electionMutex.Lock() // Ensure only one election process runs at a time
	defer n.electionMutex.Unlock()

	n.mu.Lock()
	if !n.isRunning { // Check if this node is running
		n.mu.Unlock()
		return
	}
	n.mu.Unlock()

	fmt.Printf("Node %d is starting an election.\n", n.id)

	// Assume the election is successful initially
	success := true

	for _, node := range n.nodes {
		if node.id <= n.id || !node.isRunning {
			continue // No need to send election message to nodes with ID less than or equal to current node
		}

		// If a node with a higher ID responds, the current election is overthrown
		if node.respondToElection(n.id) {
			fmt.Printf("Node %d responded with NO to %d's election request.\n", node.id, n.id)
			success = false
			break
		}
	}

	if success {
		n.mu.Lock()
		defer n.mu.Unlock()
		n.isCoordinator = true
		fmt.Printf("Node %d has become the coordinator.\n", n.id)
		fmt.Printf("Node %d start announcement\n", n.id)
		for _, node := range n.nodes {
			if n.failDuringAnnouncement {
				fmt.Printf("Node %d failed during the announcement.\n", n.id)
				n.isRunning = false
				break
			}
			node.receiveAnnouncement(n)
		}
	}
}

func (n *Node) receiveAnnouncement(newCoordinator *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.isRunning {
		fmt.Printf("Node %d did not receive the announcement", n.id)
		return
	}
	if n.failDuringElection {
		n.isRunning = false
		fmt.Printf("Node %d failed during the announcement and hence did not receiveAnnouncement.\n", n.id)
		return
	}
	if n.id < newCoordinator.id {
		n.isCoordinator = false
		n.coordinator = newCoordinator
		fmt.Printf("Node %d has received the announcement that Node %d is the new coordinator.\n", n.id, newCoordinator.id)
	} else {
		go n.startElection()
		fmt.Printf("Node %d has received the announcement that Node %d is the new coordinator. But it disapprove of the arrangement and start new election\n", n.id, newCoordinator.id)
	}
	return
}

func (n *Node) respondToElection(callerID int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.failDuringElection {
		n.isRunning = false
		fmt.Printf("Node %d failed during the election.\n", n.id)
		return false
	}

	if !n.isRunning {
		return false // This node is not running and hence won't respond
	}

	if n.id > callerID {
		// This node has a higher ID, so respond negatively and start a new election
		go n.startElection() // Start a new election since this node has a higher ID
		return true
	}

	return false // Positive response, do not start a new election
}

func (n *Node) syncData(coordinator *Node) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !coordinator.isRunning {
		return false // This node is not running, so no synchronization will occur
	}

	n.data = coordinator.data // synchronize the data
	fmt.Printf("Node %d synchronized data with the coordinator %d.\n", n.id, coordinator.id)

	return true
}

func (n *Node) start() {
	n.startElection()
	n.isCoordinator = true
}

func (n *Node) stop() {
	n.mu.Lock()
	defer n.mu.Unlock()
	fmt.Printf("Node %d is stopping.\n", n.id)
	n.isRunning = false // This will stop the node's run loop
}

func main() {
	var nodes []*Node

	// Initialize nodes
	originalCoordinator := &Node{
		id:            5,
		data:          0,
		isCoordinator: true, // Node 0 is initially the coordinator
		isRunning:     true, // All nodes start as running
	}
	for i := 0; i < 5; i++ { // Create 5 nodes
		node := &Node{
			id:            i,
			data:          0,
			isCoordinator: false, // Node 0 is initially the coordinator
			coordinator:   originalCoordinator,
			isRunning:     true, // All nodes start as running
		}
		nodes = append(nodes, node)
	}
	nodes = append(nodes, originalCoordinator)

	// Let each node know about all other nodes
	for _, node := range nodes {
		node.nodes = nodes
		go node.run() // Start each node's routine
	}

	var caseNumber int
	flag.IntVar(&caseNumber, "case", 1, "Case number to demonstrate")
	flag.Parse()

	// Wait for the system to stabilize
	fmt.Println("Waiting for the system to stabilize...")
	time.Sleep(2 * time.Second)
	switch caseNumber {
	case 1:
		fmt.Println("\n--- Case 1: Best case scenario ---")
		nodeToFail := nodes[5] // Let's assume the coordinator (node 5) fails
		nodeToFail.stop()      // Coordinator fails
		fmt.Println("Coordinator (Node 5) has failed.")

		time.Sleep(2 * time.Second) // Wait for nodes to detect the coordinator's failure

		// In the best case scenario, the node with the highest ID among the active ones initiates the election.
		// Assuming that nodes are stored in ascending order by ID and node 4 is the highest ID after node 5.
		highestIDNode := nodes[4]
		highestIDNode.startElection() // Node with the highest ID initiates the election

		fmt.Println("Node with the highest ID (Node 4) is initiating the election.")

		time.Sleep(2 * time.Second) // Wait for election to conclude and for the new coordinator to be established

		fmt.Printf("Election complete. Node %d is now the new coordinator.\n", highestIDNode.id)
		os.Exit(0)
	case 2:
		// CASE 2: Worst case scenario: Node with the lowest ID initiates the election
		fmt.Println("\n--- Case 2: Worst case scenario ---")
		nodeToFail := nodes[5] // Let's assume node 0 fails
		nodeToFail.stop()      // Node 0 fails
		fmt.Println("Coordinator has failed.")
		time.Sleep(2 * time.Second) // Wait for nodes to detect the coordinator's failure
		nodes[0].startElection()    // Node with the lowest ID initiates the election
		time.Sleep(2 * time.Second) // Wait for election to conclude

	case 3:
		// CASE 3: Node failure during election
		fmt.Println("\n--- Case 3: Node failure during election ---")

		// Trigger an election by stopping the coordinator
		fmt.Println("Coordinator (Node 5) is going down.")
		originalCoordinator := nodes[5]

		time.Sleep(2 * time.Second) // Allow time for nodes to recognize the coordinator's absence

		// 3a: The newly elected coordinator fails during the announcement
		fmt.Println("\n--- Case 3a: Newly elected coordinator fails during announcement ---")
		// Assuming node 4 is likely to be the new coordinator, we'll simulate its failure during the announcement.
		failingCoordinator := nodes[4]
		failingCoordinator.failDuringAnnouncement = true
		originalCoordinator.stop() // Coordinator fails

		time.Sleep(5 * time.Second) // Allow time for election and potential failure during announcement
	case 4:
		fmt.Println("\n--- Case 3: Node failure during election ---")

		// Trigger an election by stopping the coordinator
		fmt.Println("Coordinator (Node 5) is going down.")
		originalCoordinator := nodes[5]

		time.Sleep(2 * time.Second) // Allow time for nodes to recognize the coordinator's absence

		// 3b: A non-coordinator node fails during the election
		fmt.Println("\n--- Case 3b: Non-coordinator node fails during the election ---")
		// Assuming node 3 will fail during the election process
		failingNode := nodes[3]
		failingNode.failDuringElection = true

		// node failure necessitates a new election
		fmt.Println("Node 5 is going down, triggering a new election.")
		originalCoordinator.stop()

		time.Sleep(5 * time.Second) // Allow time for election and potential node failure during the process
	case 5:
		// CASE 4: Multiple simultaneous elections
		fmt.Println("\n--- Case 4: Multiple simultaneous elections ---")
		for _, node := range nodes {
			if node.isRunning && !node.isCoordinator { // Non-coordinator nodes that are running initiate election
				go node.startElection() // Start elections simultaneously
			}
		}
		time.Sleep(2 * time.Second) // Wait for elections to conclude

	case 6:

		// CASE 5: Node leaving silently
		fmt.Println("\n--- Case 5.1: Node leaving silently ---")
		leavingNode := nodes[2] // Let's assume node 2 leaves the network
		leavingNode.stop()      // Node 2 leaves silently
		fmt.Println("Node has left the network silently.")
		time.Sleep(2 * time.Second) // Wait for the system to react
	case 7:

		fmt.Println("\n--- Case 5.2: Coordiantor Node leaving silently ---")
		leavingNode := nodes[0] // Let's assume node 2 leaves the network
		leavingNode.stop()      // Node 2 leaves silently
		fmt.Println("Coordinator Node has left the network silently.")
		time.Sleep(2 * time.Second) // Wait for the system to react

	}
}
