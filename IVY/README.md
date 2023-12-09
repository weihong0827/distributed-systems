# Question 1: First implement the Ivy architecture discussed in the class.

## Run the project
Use the following command to run the project
```bash
go run main.go
```

## HTTP Server
To interact with the service, I implemented a simple HTTP server. The server will listen to the port 8080 and accept the following requests:
## Write
```bash
curl "localhost:8080/write?node=1&page=pg2&content=abc"
```
This command initiate a write request to the node 1, to write to page `pg2` with content `abc`. 
You can observe the output on the server side:
```bash
2023/11/28 10:43:11 node 1: initiate write for page:pg2 with content abc
2023/11/28 10:43:11 manager: write request for page pg2
2023/11/28 10:43:11 append to request
2023/11/28 10:43:11 requests: [page:"pg2" source:1]
2023/11/28 10:43:11 manager: serving write request for page pg2
2023/11/28 10:43:11 node 1: send write confirmation for page pg2
2023/11/28 10:43:11 manager: write confirmation for page pg2
2023/11/28 10:43:11 Current Page on CM: &{abc 1 []}
2023/11/28 10:43:11 Current Page on node: &{abc 1}
```
You can see that node 1 initiate the write request, and the manager received the request and append it to the request queue. Then the manager serve the request and send the confirmation back to the node 1. 

### different write coniditions
1. If you perform a write command on the same node for the same page, it already have the write access to the page and it will directly send a write confirmation to the manager. 
2. If a new page write request is sent, the manager will not send any forward request to the other nodes. It will just reply to the request node and the request node will send a write confirmation
3. If a write request is sent to the node that does not have the write access to the page, the manager will forward the request to the node that has the write access to the page. and invalidates the page in all the copy set.


## Read
```bash
curl "localhost:8080/read?node=2&page=pg2"
```
This command initiate a read request to the node 2, to read page `pg2`. In which node 2 does not have the page
If you perform this command after the write command, you will see the following output:
```bash
2023/11/28 10:47:52 Node 2 InitRead request: page:"pg2"
2023/11/28 10:47:52 Manager receive read request: page:"pg2" source:2
2023/11/28 10:47:52 Node 1 forwarding read request for page: pg2 to Node localhost:5002
2023/11/28 10:47:52 Node 2 received content of request: page:"pg2" content:"abc"
2023/11/28 10:47:52 Node 2 received page: {abc pg2}
2023/11/28 10:47:52 Read confirmation received at manager for request: page:"pg2" source:2 Read Complete
```

The current state on the CM is:
```bash
2023/11/28 10:51:56 Node 2 synchronizing state get map[pg2:copySet:2 content:"abc" owner:1]
```
we can tell that for page 2 it maintains the copy set as 2, and the owner is 1. 

### Other conditions
1. If the node has at least the read access to the page, it will directly return the page content to the request node.
2, If the node does not have the read access to the page, it will forward the request to the node that has the read access to the page. 
