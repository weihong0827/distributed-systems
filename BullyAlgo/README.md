# Bully Algorithm Implementation in Go

This README outlines the behavior of the Bully Algorithm implemented in Go. The implementation models a distributed system where nodes can initiate an election to decide on a coordinator. The coordinator synchronizes data across the nodes. The system is fault-tolerant, handling node failures during elections and node operations.

## Expected Behaviors

### Case 1: Best Case Scenario

- If the coordinator fails (e.g., Node 5), the node with the highest ID among the active ones realizes there's no coordinator and initiates an election.
- In the best-case scenario, the system quickly stabilizes with the highest ID node becoming the new coordinator without contention or further node failures.

### Case 2: Worst Case Scenario

- The coordinator fails.
- The node with the lowest ID realizes there's no coordinator and initiates the election.
- The election message traverses through most (or all) higher ID nodes, causing a delay in election completion.
- The system eventually stabilizes with the highest ID node becoming the coordinator, after potentially multiple rounds of elections.

### Case 3: Node Failure During Election

#### Case 3a: Newly Elected Coordinator Fails During Announcement

- During the announcement phase, the newly elected coordinator fails.
- This failure triggers another round of elections, causing more delays before the system stabilizes with a new coordinator.

#### Case 3b: Non-coordinator Node Fails During Election

- A non-coordinator node fails during the election process (not necessarily the one initiating the election).
- This failure can potentially lead to a new round of elections if the failed node was a candidate for becoming the coordinator.

### Case 4: Multiple Simultaneous Elections

- Multiple non-coordinator nodes initiate the election process simultaneously.
- This scenario leads to contention and potential message floods but eventually stabilizes when one node wins the election and successfully announces itself as the coordinator.

### Case 5: Node Leaving Silently

#### Case 5.1: Non-coordinator Node Leaving Silently

- A non-coordinator node leaves the network without notifying other nodes.
- The system continues operating with the remaining nodes, potentially with less efficiency until a new election happens for another reason.

#### Case 5.2: Coordinator Node Leaving Silently

- The coordinator node leaves the network without notifying other nodes.
- This absence is eventually detected by other nodes, triggering a new election process.
- The system stabilizes after the election with a new coordinator.

## Running the Code

Use the `Makefile` command to run the code

```sh
# Case 1: Best Case Scenario
make case1

# Case 2: Worst Case Scenario
make case2

# Case 3: Node Failure During Election
# case 3a - Newly Elected Coordinator Fails During Announcement
make case3.a
# case 3b - Non-coordinator Node Fails During Election
make case3.b

# Case 4: Multiple Simultaneous Elections
make case4

# Case 5: Node Leaving Silently
# case 5.1 - Non-coordinator Node Leaving Silently
make case5.a
# case 5.2 - Coordinator Node Leaving Silently
make case5.b
```

## Example Output Explanation: Case 3a - Node Failure During Election

This example demonstrates a scenario under Case 3, specifically Case 3a, where the newly elected coordinator fails during its announcement.

### Output:
```
base) ➜  q2_2 make case3.a
go run main.go -case=3
Waiting for the system to stabilize...

--- Case 3: Node failure during election ---
Coordinator (Node 5) is going down.
Node 4 synchronized data with the coordinator 5.
Node 2 synchronized data with the coordinator 5.

--- Case 3a: Newly elected coordinator fails during announcement ---
Node 5 is stopping.
Node 1 detected no coordinator
Node 1 is starting an election.
Node 2 responded with NO to 1's election request.
Node 2 is starting an election.
Node 3 responded with NO to 2's election request.
Node 3 is starting an election.
Node 4 responded with NO to 3's election request.
Node 4 is starting an election.
Node 4 has become the coordinator.
Node 4 start announcement
Node 4 failed during the announc
````


### Explanation:

1. **Coordinator (Node 5) is going down:** The system begins by taking down the current coordinator (Node 5), initiating a scenario where an election is needed.
     this critical phase, simulating a failure during the announcement.

5. **Node 3 detects no coordinator and starts a new election:** After Node 4's failure, Node 3 no longer recognizes any coordinator, prompting it to start a new election process.

6. **Node 3 becomes the new coordinator:** This time, Node 3 successfully completes the election and announces itself as the new coordinator without any failure.

7. **Other nodes acknowledge Node 3 as coordinator:** The remaining nodes receive the announcement and recognize Node 3 as the new coordinator. They then synchronize their data with Node 3, as seen in the final messages.

This case demonstrates the resilience of the bully algorithm, where the system eventually finds a new coordinator even after successive failures. However, it also highlights a potential issue where multiple failures can lead to a longer time without a coordinator, leading to the need for multiple elections.

2. **Nodes detect no coordinator:** The nodes start to realize the absence of the coordinator. Node 1 detects this first and begins an election.

3. **Election process commences:** Nodes start elections in increasing order of their IDs, as they reject candidates with IDs lower than theirs. This is seen where Node 2 rejects 1, Node 3 rejects 2, and so forth.

4. **Node 4 becomes the coordinator but fails immediately:** Node 4 successfully concludes its election and attempts to announce itself as the new coordinator. However, it was predetermined that Node 4 would fail during 
