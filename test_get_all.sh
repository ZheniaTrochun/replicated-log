# read from master
curl http://localhost:8080/get-all
# read from first sentinel
curl http://localhost:8081/get-all
# read from second sentinel
curl http://localhost:8082/get-all
