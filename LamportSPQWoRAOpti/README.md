# Lamport's Shared Priority Queue Implementation

## Overview

This Go program implements Lamport's shared priority queue for managing access to a critical section in a distributed system. It does not use Ricart and Agrawala's optimization. The algorithm uses message passing for synchronization between processes to ensure mutual exclusion without a centralized controller.

## Main Logic

The program defines a `LamportQueue` struct to represent each process's local view of the system. It includes:
- A mutex for synchronization.
- A priority queue to order requests based on Lamport timestamps.
- Logical clocks to assign timestamps.
- Channels for communication between processes.

Each `LamportQueue` runs a `Listen` routine to handle incoming messages and can broadcast requests (`RequestCS`) and release notifications (`ReleaseCS`) to other processes.

## Print Statements Explained

1. **Request Broadcast**: When a process wants to enter the critical section (CS), it broadcasts a request to all other processes.
   - Example: "Process 1 broadcasting its request at clock 2."

2. **Replies Received**: The process counts the number of replies received from other processes, indicating that it can proceed when all are collected.
   - Example: "Process 1 received reply. Total replies: 2"

3. **Entering CS**: Before entering the CS, the process checks if its request is at the head of the queue.
   - Example: "Node 1 Entering Critical Section"

4. **Release Broadcast**: Upon leaving the CS, the process removes its request from the queue and broadcasts a release message.
   - Example: "Process 1 releasing its request."

5. **Receive Request**: When a process receives a request from another process, it updates its clock and queue.
   - Example: "Process 2 received request from Process 1 at local clock 3."

6. **Send Reply**: If the process is not in the CS and does not have a prior request, it sends a reply to the requesting process.
   - Example: "Process 2 sending reply to Process 1"

7. **In CS**: Indicates that a process is currently in the critical section.
   - Example: "Process 1 in critical section"

Each print statement represents a significant event in the algorithm's execution, providing insight into the process's state and actions.

## Additional Functions

- `NewLamportQueue`: Constructor for initializing a new `LamportQueue`.
- `Less`: Helper function to compare requests based on timestamps and process IDs.
- `max`: Utility function to find the maximum of two integers.

## Main Function

The `main` function initializes channels and processes, starts listening routines, and simulates process requests to enter and exit the critical section.


# Example output walk through
```
(base) ➜  LamportSPQWoRAOpti git:(main) ✗ go run main.go
Process 3 broadcasting its request at clock 1.
Process 1 received request from Process 3 at local clock 0.
Process 1 sending reply to Process 3
Process 3 received reply. Total replies: 1
Process 2 received request from Process 3 at local clock 0.
Process 2 sending reply to Process 3
Process 3 received reply. Total replies: 2
Node 3 Entering Critical Section
Process 3 in critical section
Process 1 broadcasting its request at clock 3.
Process 3 received request from Process 1 at local clock 1.
Process 2 broadcasting its request at clock 3.
Process 1 received request from Process 2 at local clock 3.
Process 2 received request from Process 1 at local clock 3.
Process 2 sending reply to Process 1
Process 1 received reply. Total replies: 1
Process 3 received request from Process 2 at local clock 4.
Process 3 releasing its request.
Process 2 received reply. Total replies: 1
Process 1 received reply. Total replies: 2
Node 1 Entering Critical Section
Process 1 in critical section
Process 1 releasing its request.
Process 2 received reply. Total replies: 2
Node 2 Entering Critical Section
Process 2 in critical section
Process 2 releasing its request.
```

The above code is run with `nClients=3` for demonstration purpose. You can change the number of clients using the variable. The above code simulates 3 node requesting to enter CS at the same time

1. Each processes starts to request to enter critical section 
2. The process `broadcast` it request to all other nodes
3. After receiving request, if there is no pending request, simply `reply`
4. once all replies from the other nodes are received, current node enter CS
5. If a node that receive a request and it currently has a higher priority request pending. It will hold on to the reply
6. After node 3 finish its CS, it sends a release message back to node 1 and 2. Upon receiving a `release` message, node 1 and 2 both pop the head of the queue and take node 3 as replied to their pending requests
7. since node 1 has all the replies from node 2 and 3 while holding on to the reply to node 2. node 1 will enter CS but node 2 cannot since node 1's request has a higher priority than node 2

