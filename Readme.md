# Replicated log
  
**This is a toy project for trying out Golang and learning some distributed systems concepts.**   
  
## Idea
Master receives messages via `POST /insert` endpoint. It stores message in memory and replicates it to all sentinels.   
Consistency level of insert operation is configurable using `writeConcern` request field - number of instances that needs to persist item in order to mark request as success.  
Request format: `{"message": string, "writeConcern": int}`  
  
Sentinels expose `gRPC` endpoint that stores messages replicated from master (and should be called only by master). Request format: `"id": int32, "message": string, "timestamp": int64`
  
All sentinels and master expose `GET /get-all` endpoint that returns list of all stored messages.  

## Requirements
1. All messages on all instances should be stored in the insertion (to master) order.
2. Successful response for `/insert` endpoint should be returned only after successful replication to all instances.
3. `GET /get-all` on all instances should return the same messages in the same order (if no concurrent insertions during `get-all` call).

## Limitations
For now we assume that all requests to sentinels successful and networks always work reliably and there is no additional delays.

## Example of work
![Example](doc/example_of_work.png)
*Note: Logs may a bit out of order between containers.*

## How to run locally

```bash
docker-compose build
docker-compose up
```

These commands will start master on port `8080` and two sentinels on ports `8081` and `8082`.  
Operating locally:

```bash
./test_1.sh &
sleep 3
./test_2.sh &
sleep 17
./test_get_all.sh
```

Cleanup:

```bash
docker-compose down
```

## Configs
All configs are passed as env variables:  
- ROLE: "master" or "sentinel" - role of instance in a cluster - REQUIRED 
- PORT: port for HTTP server - OPTIONAL - 8080 by default
- GRPC_PORT: port for gRPC server - only for sentinel applications - OPTIONAL - 9090 by default
- SENTINELS: coma-separated list of sentinel instances in format \<host:grpc_port\> - REQUIRED for master

## Algorithm for insertion
1. Store new message locally on master
   1. Acquire lock
   2. Generate new id as last stored id + 1
   3. Insert new item by new id
   4. Release lock
2. For each sentinel in parallel
   1. Send `gRPC` request to sentinel with message and id
   2. Sentinel inserts item locally by id without locking
   3. Sentinel responds with success
   4. If response successful - send success to result channel
3. Wait for `writeConcern - 1` responses from sentinels
4. Return success
