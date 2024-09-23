# Replicated log
  
**This is a toy project for trying out Golang and learning some distributed systems concepts.**   
  
### Idea
Master receives messages via `POST /insert` endpoint. It stores message in memory and replicates it to all sentinels.  
Sentinels expose `POST /replicate-item` endpoint that stores messages replicated from master (and should be called only by master).  
All sentinels and master expose `GET /get-all` endpoint that returns list of all stored messages.  

### Requirements
1. All messages on all instances should be stored in the insertion (to master) order.
2. Successful response for `/insert` endpoint should be returned only after successful replication to all instances.
3. `GET /get-all` on all instances should return the same messages in the same order (if no concurrent insertions during `get-all` call).

### Example of work

![Example](doc/example_of_work.png)

### How to run locally

```bash
docker-compose build
docker-compose up
```

These commands will start master on port `8080` and two sentinels on ports `8081` and `8082`.  
Operating locally:

```bash
# insert several messages
curl -X POST http://localhost:8080/insert -d '{"message": "1"}'
curl -X POST http://localhost:8080/insert -d '{"message": "2"}'
curl -X POST http://localhost:8080/insert -d '{"message": "3"}'

# read from master
curl http://localhost:8080/get-all
# read from first sentinel
curl http://localhost:8081/get-all
# read from second sentinel
curl http://localhost:8082/get-all
```
