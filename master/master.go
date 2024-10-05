package master

import (
	"log/slog"
	"replicated-log/repository"
	"replicated-log/sentinel"
)

type LogMaster struct {
	sentinels []*sentinel.SentinelClient
}

var service *LogMaster

func InitLogMasterService(sentinelAddresses []string) {
	sentinels := make([]*sentinel.SentinelClient, len(sentinelAddresses))

	for i, address := range sentinelAddresses {
		sentinels[i] = sentinel.NewSentinelClient(i, address)
	}

	service = &LogMaster{
		sentinels: sentinels,
	}
}

func storeMessage(msg string) (int, error) {
	item := repository.Insert(msg)

	resChannel := make(chan int, len(service.sentinels))
	errChannel := make(chan error, len(service.sentinels))

	for _, sentinelClient := range service.sentinels {
		go sentinelClient.ReplicateItem(item, resChannel, errChannel)
	}

	for range service.sentinels {
		select {
		case replicaId := <-resChannel:
			slog.Info("Replica updated", "replica_id", replicaId, "item_id", item.Id)
		case err := <-errChannel:
			return -1, err
		}
	}

	return item.Id, nil
}
