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

func storeMessage(msg string, writeConcern int) (int, error) {
	item := repository.Insert(msg)

	resChannel := make(chan int, len(service.sentinels))
	errChannel := make(chan error, len(service.sentinels))

	for _, sentinelClient := range service.sentinels {
		go sentinelClient.ReplicateItem(item, resChannel, errChannel)
	}

	wc := sanitizeWriteConcern(writeConcern)

	for range wc - 1 {
		select {
		case replicaId := <-resChannel:
			slog.Info("Replica updated", "replica_id", replicaId, "item_id", item.Id)
		case err := <-errChannel:
			return -1, err
		}
	}

	return item.Id, nil
}

func sanitizeWriteConcern(writeConcern int) int {
	if writeConcern < 1 || writeConcern > len(service.sentinels) {
		return len(service.sentinels)
	}

	return writeConcern
}
