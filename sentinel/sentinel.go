package sentinel

import "replicated-log/persistence"

type Sentinel struct {
	store *persistence.DataStore
}

func NewSentinel() *Sentinel {
	return &Sentinel{persistence.NewStore()}
}

func (s *Sentinel) Sync(id int, value string, timestamp int64) bool {
	item := persistence.Item{id, value, timestamp}

	return s.store.InsertById(item)
}

func (s *Sentinel) GetAll() []string {
	messages := s.store.GetAll()

	res := make([]string, len(messages))

	for i, message := range messages {
		res[i] = message.Value
	}

	return res
}
