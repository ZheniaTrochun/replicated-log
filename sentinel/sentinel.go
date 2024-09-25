package sentinel

import "replicated-log/repository"

func sync(id int, value string, timestamp int64) bool {
	item := repository.Item{id, value, timestamp}

	return repository.InsertById(item)
}
