FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /replicated-log

ENV GIN_MODE=release

EXPOSE 8080

CMD ["/replicated-log"]
