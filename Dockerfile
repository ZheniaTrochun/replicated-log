FROM golang:1.23.1

WORKDIR /app

RUN apt-get update && apt-get install --no-install-recommends --assume-yes protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN export PATH="$PATH:$(go env GOPATH)/bin"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sentinel/sentinel.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /replicated-log

ENV GIN_MODE=release

EXPOSE 8080
EXPOSE 9090

CMD ["/replicated-log"]
