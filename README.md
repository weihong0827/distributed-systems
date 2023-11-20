# Programming assignment
This is for SUTD 50.041 Distributed System and Computing course programming assignment 1. In this assignment, we will be implementing some basic but core distributed concepts. This entire assignment is done in golang.

# Assignment 1
## Clocks

This section implements a simple distributed client server architecture, simulated using goroutine. We will be implementing Lamport's logical clock to estamish the ordering of events in a distributed system. The client will send a message to the server, and the server will send a message back to the client. The client will then print out the message it received from the server, along with the timestamp of the message. The server will also print out the message it received from the client, along with the timestamp of the message. 

We will also be implementing Vector clock to for any causality violation.

Both implementations will be running with at least 10 clients, but the number of clients can be changed base on the user input.

### Where to find the code
The code and implementations for this section can be found in the `Clocks` folder.
The explanation is in the [README.md](./Clocks/README.md) file in the `Clocks` folder.

## Bully Algorithm

This section implements a simple bully algorithm to run election for a distributed system and take care of node failures. The system will have a coordinator that will be responsible for synchronizing data across the nodes. The system is fault-tolerant, handling node failures during elections and node operations.

### Where to find the code
The code and implementations for this section can be found in the `BullyAlgo` folder.
The explanation is in the [README.md](./BullyAlgo/README.md) file in the `BullyAlgo` folder.


# Assignment 2
There are 4 different parts to this Assignment focusing on distributed mutex implementations
1. Implement Lamport's shared priority queue without Ricart and Agrawala's optimization
2. Implement Lamport's shared priority queue with Ricart and Agrawala's optimization
3. Implement voting protocol with deadlock avoidance (vote rescinding)

4. Compare the performance of three implementations

# Performance comparision

| Time Taken for Protocol (ms)      | 1    | 2    | 3    | 4    | 5    | 6    | 7    | 8    | 9    | 10   |
|-----------------------------------|------|------|------|------|------|------|------|------|------|------|
| Lamport Shared Priority Queue     | 0.62 | 1.30 | 1.45 | 1.50 | 2.75 | 2.90 | 3.30 | 3.50 | 4.15 | 5.00 |
| Ricart & Agrawala                 | 0.35 | 0.40 | 0.42 | 0.44 | 0.78 | 0.83 | 1.30 | 1.40 | 1.68 | 1.77 |
| Voting Protocol                   | 0.12 | 0.14 | 0.40 | 0.55 | 0.70 | 0.82 | 0.95 | 1.25 | 1.43 | 1.70 |




