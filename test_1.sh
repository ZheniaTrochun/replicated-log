curl -X POST http://localhost:8080/insert -d '{"message": "1", "writeConcern": 3}'
curl -X POST http://localhost:8080/insert -d '{"message": "2", "writeConcern": 2}'
curl -X POST http://localhost:8080/insert -d '{"message": "3", "writeConcern": 3}'

./test_get_all.sh
