curl -X POST http://localhost:8080/insert -d '{"message": "4", "consistency": 1}'
curl -X POST http://localhost:8080/insert -d '{"message": "5", "consistency": 1}'
curl -X POST http://localhost:8080/insert -d '{"message": "6", "consistency": 1}'

# read from master
curl http://localhost:8080/get-all
# read from first sentinel
curl http://localhost:8081/get-all
# read from second sentinel
curl http://localhost:8082/get-all
