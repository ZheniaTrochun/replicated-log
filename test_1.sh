curl -X POST http://localhost:8080/insert -d '{"message": "1", "consistency": 3}'
curl -X POST http://localhost:8080/insert -d '{"message": "2", "consistency": 2}'
curl -X POST http://localhost:8080/insert -d '{"message": "3", "consistency": 3}'

# read from master
curl http://localhost:8080/get-all
# read from first sentinel
curl http://localhost:8081/get-all
# read from second sentinel
curl http://localhost:8082/get-all