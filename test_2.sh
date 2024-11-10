curl -X POST http://localhost:8080/insert -d '{"message": "4", "writeConcern": 1}'
curl -X POST http://localhost:8080/insert -d '{"message": "5", "writeConcern": 1}'
curl -X POST http://localhost:8080/insert -d '{"message": "6", "writeConcern": 1}'

./test_get_all.sh
