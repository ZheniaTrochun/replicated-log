FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#COPY *.go .
#COPY persistence/*.go ./persisitence/
#COPY master/*.go ./master/.
#COPY sentinel/*.go ./sentinel/.

RUN CGO_ENABLED=0 GOOS=linux go build -o /replicated-log

EXPOSE 8080

CMD ["/replicated-log"]