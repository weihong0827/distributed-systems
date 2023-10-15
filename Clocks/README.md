# Distributed Messaging Simulation in Go

This project simulates a simple distributed system with a server and multiple clients communicating asynchronously. It showcases the concepts of logical and vector clocks in the context of message ordering and causal relationships.

## Overview

The system comprises clients (`Client`) and a server (`Server`) that exchange messages. Each `Client` and `Server` maintains its own logical clock or vector clock, used to timestamp messages and determine the causality and order of events.

### Files

- `client.go`: Contains the `Client` struct and methods for client operations, including sending messages to and receiving messages from the server.
- `server.go`: Contains the `Server` struct and methods for server operations, including broadcasting messages to clients and handling received messages.

## Key Components

### Client

Each client periodically sends a message to the server and listens for incoming messages. The client can operate using a simple logical clock or a more complex vector clock, based on the `useVectorClock` flag.

#### Main Functions

- `Send()`: Periodically creates and sends a message to the server. It increments its own clock (logical or vector) with each sent message.
- `Receive()`: Listens for incoming messages from the server, updates its own clock based on the received message, and checks for causal violations.

### Server

The server listens for messages from clients. Upon receiving a message, it can either forward it to all other clients or simulate a message drop. The server maintains its own clock and updates it based on incoming messages or messages it decides to send/drop.

#### Main Functions

- `Start()`: Begins the server's listen loop. On each loop iteration, the server receives a message, decides whether to drop or forward it, updates its own clock, and sends the message to all clients except the sender.

## Clocks and Causality

This simulation uses logical and vector clocks to establish a partial ordering of events and detect causality violations.

- **Logical Clock**: A simple counter incremented with each new event. It doesn't provide enough information to distinguish between concurrent events.
- **Vector Clock**: An array of counters, one for each client and the server. It provides more detailed information about the system's state, allowing it to distinguish concurrent events and detect causal violations.

## Scenarios and Expected Outputs

### Case 1 and 2: Messages with Logical Clocks

When a client sends a message, it includes its current clock value. The server, upon receiving the message, may decide to drop it (simulating network loss) or broadcast it to other clients. When the message is not dropped, it's timestamped with the server's current clock and forwarded to all clients except the sender.

**Client's Message Sending**

- Output: `"Hello from client %d at time %d", c.id, c.clock` or `"Hello from client %d at time %d", c.id, c.GetVectorClock()[c.id]`
  - Explanation: This prints a message that the client is sending, tagged with the client's ID and the current time from the logical or vector clock.

**Server's Message Receiving and Dropping**

- Output: `"Server received: %s ,at server clock %d\n", msg.Msg, s.clock` or `"Server received: %s ,at server vector clock %v\n", msg.Msg, s.vectorClock`
  - Explanation: Indicates that the server has received a message and shows the message content along with the server's current logical or vector clock value.

- Output: `"Dropped:%s at server clock %d \n", msg.Msg, s.clock` or `"Dropped:%s at server vector clock %v \n", msg.Msg, s.vectorClock`
  - Explanation: Indicates that the server has decided to simulate a network loss by dropping a message. It shows the dropped message's content and the server's current clock.

**Server's Message Forwarding**

- Output: `"Server sent:%s to client %d at server clock %d \n", updatedMsg.Msg, client.id, s.clock` or `"Server sent:%s to client %d at server vector clock %v \n", updatedMsg.Msg, client.id, s.vectorClock`
  - Explanation: This prints the content of the message that the server is forwarding to clients, excluding the sender. It shows the message, the ID of the receiving client, and the server's current logical or vector clock value.

#### Final Sequence printing
At the end of 5 seconds, all clients and server will stop execution and each client will print out the message that they received in order

### Case 3: Vector Clock

Causal violations occur when a client receives a message with a clock indicating that it should have received some prior message, but it hasn't. This is detected by comparing the received message's vector clock with the client's own vector clock.

**Causal Violation Detection**

- Output: `"Client %d found causal violation with msg clock %v and client clock %v \n", c.id, msg.GetVectorClock(), c.GetVectorClock()` or `"Server found causal violation with msg clock %v and server clock %v \n", msg.GetVectorClock, s.GetVectorClock`
  - Explanation: These statements print when the system detects a causal violation, meaning a message was received that indicates a preceding message was missed. It shows the clocks of both the received message and the recipient, highlighting the discrepancy.

**Regular Message Receiving**

- Output: `"Client %d received:%s with msg clock %v and client clock %v \n", c.id, msg.Msg, msg.GetVectorClock(), c.GetVectorClock()` or `"Client %d received: %s ,at client clock %d\n", c.id, msg.Msg, c.GetClock()`
  - Explanation: Prints the message a client has received and the corresponding clocks without any causal violation. It's an indicator of normal operation, showing the message content, the client's ID, and the relevant clock values.

## Running the Code

To run the code, simpily run `make run<question num>` in this case it can be `make run1` or `make run2` or `make run3`
The defult number of client is 10, but you can change it by `make run<question num> NUM_CLIENTS=<number>`

### Sample Output and Explanation
In this case i will be using `make run2 NUM_CLIENTS=3` for easy explanation

#### Sending and Receiving Messages between Clients and Server
In the following output, the server receives messages from the client and increment its own clock, then increment its own clock again and sends the message to other clients. The client receives the message and compare and increment its own clock with the clock in the message. The server's clock is always ahead of the clients' clocks, since it receives and sends messages before the clients do.
```
Server received: Hello from client 2 at time 1 ,at server clock 2
Dropped:Hello from client 2 at time 1 at server clock 3
Server received: Hello from client 0 at time 1 ,at server clock 4
Server sent:Hello from client 0 at time 1 to client 1 at server clock 5
Server sent:Hello from client 0 at time 1 to client 2 at server clock 6
Server received: Hello from client 1 at time 1 ,at server clock 7
Server sent:Hello from client 1 at time 1 to client 0 at server clock 8
Server sent:Hello from client 1 at time 1 to client 2 at server clock 9
Client 2 received: Hello from client 0 at time 1 ,at cleint clock 7
Client 2 received: Hello from client 1 at time 1 ,at cleint clock 10
Client 1 received: Hello from client 0 at time 1 ,at cleint clock 6
Client 0 received: Hello from client 1 at time 1 ,at cleint clock 9
```

#### Question 2 finial Sequence
Below shows a sample output of the final sequence of messages received by each client, in order. The messages are tagged with the client's ID and the time at which it was received. The messages are printed in order of receipt, but the time values may not be in order due to the asynchronous nature of the system.
```
Client 0 messages:
At Client Clock 8 received: Hello from client 1 at time 1
At Client Clock 16 received: Hello from client 2 at time 11
At Client Clock 27 received: Hello from client 1 at time 20
Client 1 messages:
At Client Clock 5 received: Hello from client 0 at time 1
At Client Clock 17 received: Hello from client 2 at time 11
At Client Clock 32 received: Hello from client 0 at time 19
Client 2 messages:
At Client Clock 6 received: Hello from client 0 at time 1
At Client Clock 9 received: Hello from client 1 at time 1
At Client Clock 28 received: Hello from client 1 at time 20
At Client Clock 33 received: Hello from client 0 at time 19
```

#### VectorClock Output
Running `make run3 NUM_CLIENTS=3` will show the output of the vector clock
```
Server received: Hello from client 1 at time 1 ,at server vector clock [0 1 0 1]
Dropped:Hello from client 1 at time 1 at server vector clock [0 1 0 2]
Server received: Hello from client 2 at time 1 ,at server vector clock [0 1 1 3]
Server sent:Hello from client 2 at time 1 to client 0 at server vector clock [0 1 1 4]
Server sent:Hello from client 2 at time 1 to client 1 at server vector clock [0 1 1 5]
Server received: Hello from client 0 at time 1 ,at server vector clock [1 1 1 6]
Server sent:Hello from client 0 at time 1 to client 1 at server vector clock [1 1 1 7]
Server sent:Hello from client 0 at time 1 to client 2 at server vector clock [1 1 1 8]
Client 2 received:Hello from client 0 at time 1 with msg clock [1 1 1 8] and client clock [1 1 2 8]
Client 1 received:Hello from client 2 at time 1 with msg clock [0 1 1 5] and client clock [0 2 1 5]
Client 1 received:Hello from client 0 at time 1 with msg clock [1 1 1 7] and client clock [1 3 1 7]
Client 0 received:Hello from client 2 at time 1 with msg clock [0 1 1 4] and client clock [2 1 1 4]
```
same idea as before, but this time the clock is a vector clock

### Causal Violation
I have a function written to detect causality violation, but it is never triggered for my case, i will attach the function below for reference
```go
func CheckCausalViolation(entity Clock, message Message) bool {
	// Check for at least one grater than
	atLeastOneGreater := false
	for i, val := range entity.GetVectorClock() {
		if val < message.GetVectorClock()[i] {
			atLeastOneGreater = true
		}
		if val > message.GetVectorClock()[i] {
			return false
		}
	}
	return atLeastOneGreater
}
```
