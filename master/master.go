package master

import (
	"github.com/apex/log"
	"replicated-log/persistence"
	"replicated-log/sentinel"
)

type LogMaster struct {
	store     *persistence.DataStore
	sentinels []*sentinel.SentinelClient
}

func NewLogMaster(sentinelAddresses []string) *LogMaster {
	store := persistence.NewStore()

	sentinels := make([]*sentinel.SentinelClient, len(sentinelAddresses))

	for i, address := range sentinelAddresses {
		sentinels[i] = sentinel.NewSentinelClient(address)
	}

	return &LogMaster{
		store:     store,
		sentinels: sentinels,
	}
}

func (m *LogMaster) StoreMessage(msg string) (int, error) {
	item := m.store.Insert(msg)

	resChannel := make(chan int, len(m.sentinels))
	errChannel := make(chan error, len(m.sentinels))

	for _, sentinelClient := range m.sentinels {
		go sentinelClient.SyncItem(item, resChannel, errChannel)
	}

	for range m.sentinels {
		select {
		case <-resChannel:
			log.Info("Replication finished")
		case err := <-errChannel:
			return -1, err
		}
	}

	return len(m.sentinels), nil
}

func (m *LogMaster) GetAll() []string {
	messages := m.store.GetAll()

	res := make([]string, len(messages))

	for i, message := range messages {
		res[i] = message.Value
	}

	return res
}
