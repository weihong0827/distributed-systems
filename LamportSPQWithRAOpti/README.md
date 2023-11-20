# Lamport's Shared Priority Queue with Ricart and Agrawala’s Optimization

## Overview

This Go program implements Lamport's shared priority queue with Ricart and Agrawala’s optimization, a method for ensuring mutual exclusion in distributed systems. The algorithm enhances efficiency by reducing the number of messages required for a process to enter the critical section (CS).

## Main Logic

The program uses a `LamportQueue` structure that encompasses the necessary data for each process, such as the priority queue, logical clock, and channels for inter-process communication. The main logic involves the following steps:

1. Requesting access to the CS (`RequestCS`).
2. Receiving and handling requests from other processes (`ReceiveRequest`).
3. Releasing the CS and notifying other processes (`ReleaseCS`).
4. Continuously listening for incoming requests (`Listen`).

The key optimization lies in minimizing the communication overhead by only sending replies when necessary and ensuring that all non-deferred requests are replied to when the CS is released.

## Print Statements Explained

1. **Request Broadcast**: When a process wants to enter the CS, it increments its clock, broadcasts a request to all other processes, and waits for replies.
   - Example: "Process 1 broadcasting its request at clock 2."

2. **Replies Received**: Tracks the number of replies received from other processes, ensuring all replies are collected before entering the CS.
   - Example: "Process 1 received reply. Total replies: 2"

3. **Entering CS**: A process enters the CS only if its request is at the head of the queue.
   - Example: "Node 1 Entering Critical Section"

4. **Release Broadcast**: After completing its CS execution, a process releases the CS and sends replies to all pending requests.
   - Example: "Process 1 releasing its request."

5. **Receive Request**: Upon receiving a request, if the process is not in the CS or the incoming request has a higher priority, it replies immediately.
   - Example: "Process 2 received request from Process 1 at local clock 3."

6. **Send Reply**: Replies are sent to the requesting process either immediately upon receiving a request or upon releasing the CS.
   - Example: "Process 2 sending reply to Process 1"

7. **In CS**: Indicates active execution within the critical section by a process.
   - Example: "Process 1 in critical section"

These print statements provide a debug-friendly view of the system's state and interactions.

## Additional Functions

- `NewLamportQueue`: Initializes a new instance of `LamportQueue`.
- `Less`: Compares two requests to determine queue ordering.
- `max`: Returns the greater of two integers.

## Main Function

The `main` function sets up the system with multiple processes, each represented by a `LamportQueue` instance. It initiates the listening for requests and simulates the process execution flow of requesting, entering, and releasing the CS.

# Example walk through

```
(base) ➜  LamportSPQWithRAOpti git:(main) ✗ go run main.go
Process 1 broadcasting its request at clock 1.
Process 3 broadcasting its request at clock 1.
Process 1 received request from Process 3 at local clock 1.
Process 2 broadcasting its request at clock 1.
Process 3 received request from Process 1 at local clock 1.
Process 3 sending reply to Process 1
Process 3 received request from Process 2 at local clock 1.
Process 3 sending reply to Process 2
Process 2 received reply. Total replies: 1
Process 2 received request from Process 1 at local clock 1.
Process 1 received reply. Total replies: 1
Process 2 sending reply to Process 1
Process 2 received request from Process 3 at local clock 1.
Process 1 received request from Process 2 at local clock 2.
Process 1 received reply. Total replies: 2
Node 1 Entering Critical Section
Process 1 in critical section
Process 1 releasing its request.
Process 3 received reply. Total replies: 1
Process 2 received reply. Total replies: 2
Node 2 Entering Critical Section
Process 2 in critical section
Process 2 releasing its request.
Process 3 received reply. Total replies: 2
Node 3 Entering Critical Section
Process 3 in critical section
Process 3 releasing its request.
```
The logic is pretty much the same as what is in [part 1](../LamportSPQWoRAOpti/README.md)
The only difference is that the request that you replied, do not get stored in the shared pq.
During release stage, you only send release message to the pending requests stored. you can see that when node 1 exiting cs it sends release to 2 and 3. When node 2 release, it only sends to 3 but not 1. This is because node 2 knows that theres a pending request from node 3 that has not yet replied

