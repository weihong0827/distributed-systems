package main

import (
	"flag"
	"fmt"
	"sort"
	"sync"
	"time"
)

func start(numClients int, useVectorClock bool) {
	serverStopCh := make(chan bool)
	vectorClockSize := numClients + 1
	server := NewServer(vectorClockSize, useVectorClock, serverStopCh)
	stopChannels := make([]chan bool, numClients)

	for i := 0; i < numClients; i++ { // Changed to i < numClients to include all clients
		client := NewClient(i, server.broadcast, vectorClockSize, useVectorClock)
		server.register(client)
		stopChannels[i] = make(chan bool)
		client.StopCh = stopChannels[i]

		client.wg.Add(2) // Expecting two goroutines per client
		go client.Send()
		go client.Receive()

	}

	var serverWg sync.WaitGroup
	serverWg.Add(1) // Expecting one goroutine for server
	go func() {
		defer serverWg.Done() // Ensure to signal that this goroutine is done at the end
		server.Start()
	}()

	// Stop the clients after 5 seconds
	time.AfterFunc(5*time.Second, func() {
		for i := 0; i < numClients; i++ {
			close(stopChannels[i]) // Signal to stop the goroutines
		}
	})

	go func() {
		for i := 0; i < numClients; i++ {
			client := server.Clients[i]
			client.wg.Wait()
		}

		// All clients are done, send stop signal to server
		serverStopCh <- true // Send a signal instead of closing the channel
	}()

	if useVectorClock {
		serverWg.Wait()

		return
	}

	// Wait for all clients to finish processing and then print the messages
	for i := 0; i < numClients; i++ {
		client := server.Clients[i]
		client.wg.Wait()

		// Sort the messages by logical clock
		sort.Slice(client.ReceivedMsgs, func(i, j int) bool {
			return client.ReceivedMsgs[i].Clock < client.ReceivedMsgs[j].Clock
		})

		// Print the sorted messages
		fmt.Printf("Client %d messages:\n", client.id)
		for _, msg := range client.ReceivedMsgs {
			fmt.Printf("At Client Clock %d received: %s\n", msg.Clock, msg.Msg)
		}
	}
	serverWg.Wait()
}

func main() {
	question := flag.Float64("question", 1.1, "Enter the question number to run the relevant code")
	numClients := flag.Int("numClients", 5, "Enter the number of clients to run")
	useVectorClock := false
	flag.Parse()

	switch *question {
	case 1.1:
		fmt.Println("Question 1.1")
		start(*numClients, useVectorClock)

	case 1.2:
		fmt.Println("Question 1.2")
		start(*numClients, useVectorClock)
	case 1.3:
		fmt.Println("Question 1.3")
		useVectorClock = true
		start(*numClients, useVectorClock)

	default:
		fmt.Println("Question not found")
	}
}
