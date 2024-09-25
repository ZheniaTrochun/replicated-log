package base

import "replicated-log/repository"

func getAll() []string {
	messages := repository.GetAll()

	res := make([]string, len(messages))

	for i, message := range messages {
		res[i] = message.Value
	}

	return res
}
