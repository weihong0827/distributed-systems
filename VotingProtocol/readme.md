# Voting Protocol

This Go package implements a distributed mutual exclusion algorithm. The code allows multiple nodes in a distributed system to enter a critical section (CS) without conflict by using message passing and timestamps to ensure mutual exclusion.

## Overview

The code defines a distributed system where nodes can request to enter the CS, receive permission, and release the CS when they are done. It implements a variation of the Ricart-Agrawala algorithm, which is an extension of Lamport's logical clocks for distributed mutual exclusion.

## Key Components

### MessageType

`MessageType` is an enumeration representing the types of messages that can be passed between nodes. The types are `Release`, `Rescind`, and `RescindReply`.

### Message

`Message` is a struct that represents a message to be passed between nodes, containing the message type and the process ID (`pid`) of the sender.

### Request

`Request` is a struct representing a request to enter the CS, with a timestamp and the `pid` of the requester.

### Node

`Node` represents a process in the distributed system with several important fields:

- `pid`: Process ID.
- `votedFor`: The current request the node has voted for.
- `clock`: A logical clock for timestamping messages.
- `mu`: A mutex to protect critical sections.
- `queue`: A sorted queue of requests to enter the CS.
- `inCS`: A boolean indicating if the node is in the CS.
- `votedFrom`: A slice of process IDs that have voted for this node.
- `requestChan`, `messageChan`, `replyChan`: Channels for communication with other nodes.

### Methods

#### RequestCS

Allows a node to request entry to the CS.

#### ReleaseCS

Allows a node to release the CS and inform other nodes.

#### handleRequest

Handles incoming requests to enter the CS.

#### handleMessage

Handles messages related to voting and coordination between nodes.

#### Listen

Listens for incoming requests and messages and processes them.

#### sortQueue

Sorts the queue of requests based on timestamps and process IDs to ensure a total order.

## Functionality

1. When a node wants to enter the CS, it sends a request to all other nodes (`RequestCS`).
2. Each node maintains a queue of incoming requests, sorted by timestamp.
3. Nodes reply to requests based on the order in the queue.
4. When a node receives enough replies, it enters the CS.
5. Upon exiting the CS (`ReleaseCS`), the node informs all nodes that voted for it.
6. Nodes handle incoming requests and release messages, and can rescind a vote if a request with a higher priority arrives.

## Output Explanation

### RequestCS Method
"Node %d requesting CS at local clock %d\n": This statement logs when a node is requesting to enter the critical section (CS) and displays the node's ID and its local logical clock value at the time of the request.
"Node %d receive reply from %d. Total Replies:%d\n": Logs that a node has received a reply from another node. It shows the ID of the requesting node, the ID of the replying node, and the total number of replies received so far.
"Node %d Entering CS\n": Indicates that the node has received sufficient replies to enter the CS and is now entering it.
### ReleaseCS Method
"Node %d releasing vote to %d\n": Logs that a node is releasing its vote (or permission) to another node. This typically happens when the node is leaving the CS.
"Node %d exited CS\n": Indicates that the node has finished its work in the CS and has exited, allowing others to enter.
### handleRequest Method
"Node %d vote for node %d\n": When a node receives a request to enter the CS, it logs that it is voting for the requesting node by sending a reply.
"Node %d already voted, checking for priority..\n": This log statement indicates that the node has already voted for another request, and it will check if the new incoming request has a higher priority based on the timestamp.
"votedFor: %d %d, incoming %d %d\n": Shows the comparison between the currently voted-for request's ID and timestamp, and the incoming request's ID and timestamp.
"Node %d sending rescind message to node %d\n": Indicates that a node is sending a rescind message to another node, usually because a higher priority request has come in.
"Sent": A simple acknowledgment indicating that a message (typically a rescind message) has been sent.
### handleMessage Method
"Node %d currently in cs not able to rescind": When a node receives a rescind message while it is in the CS, it logs that it cannot rescind its vote at that moment.
"Node %d release vote by %d, current votes:%v\n": Logs that a node is releasing its vote due to a rescind request from another node, and shows the current list of votes.
"queue at node %d is %v\n": Shows the current state of the request queue for a node.
### Listen Method
"Receive CS request from %d at node %d\n": Logs when a node receives a request to enter the CS from another node.
"Receive message %s from %d at node %d\n": Logs when a node receives a message with the message type (e.g., Release, Rescind, RescindReply) from another node.

## Usage

Instantiate nodes with `NewNode` by passing the necessary channels for communication and a unique process ID. Start the `Listen` method on each node to begin processing messages and requests.

## Conclusion

This code provides a working example of distributed mutual exclusion using message passing. It is designed to be a learning tool for understanding how distributed systems can coordinate access to shared resources without conflicts.

# Performance comparision

| Time Taken for Protocol (ms)      | 1    | 2    | 3    | 4    | 5    | 6    | 7    | 8    | 9    | 10   |
|-----------------------------------|------|------|------|------|------|------|------|------|------|------|
| Lamport Shared Priority Queue     | 0.62 | 1.30 | 1.45 | 1.50 | 2.75 | 2.90 | 3.30 | 3.50 | 4.15 | 5.00 |
| Ricart & Agrawala                 | 0.35 | 0.40 | 0.42 | 0.44 | 0.78 | 0.83 | 1.30 | 1.40 | 1.68 | 1.77 |
| Voting Protocol                   | 0.12 | 0.14 | 0.40 | 0.55 | 0.70 | 0.82 | 0.95 | 1.25 | 1.43 | 1.70 |

