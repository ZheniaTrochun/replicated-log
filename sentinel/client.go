package sentinel

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"replicated-log/repository"
)

type SentinelClient struct {
	id      int
	address string
}

func NewSentinelClient(id int, address string) *SentinelClient {
	return &SentinelClient{id, address}
}

func (s *SentinelClient) ReplicateItem(item repository.Item, res chan int, errorChan chan error) {
	conn, err := grpc.NewClient(s.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Failed to open connection to server.", "error", err)
		errorChan <- err
		return
	}
	defer conn.Close()

	request := ReplicateRequest{Id: int32(item.Id), Message: item.Value, Timestamp: item.Timestamp}

	client := NewReplicatedLogSentinelClient(conn)

	_, err = client.Replicate(context.Background(), &request)
	if err != nil {
		slog.Error("Failed to replicate item.", "item", item, "error", err)
		errorChan <- err
		return
	}

	res <- s.id
}
