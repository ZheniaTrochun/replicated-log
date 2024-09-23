package sentinel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"replicated-log/persistence"
)

type SentinelClient struct {
	address string
}

const ReplicateEndpoint = "/replicate-item"

func NewSentinelClient(address string) *SentinelClient {
	return &SentinelClient{address}
}

func (s *SentinelClient) SyncItem(item persistence.Item, res chan int, err chan error) {
	request := Request{item.Id, item.Value, item.Timestamp}

	serializedRequest, serErr := json.Marshal(request)

	if serErr != nil {
		err <- serErr
	}

	bodyReader := bytes.NewReader([]byte(serializedRequest))

	response, httpErr := http.Post(s.address+ReplicateEndpoint, "application/json", bodyReader)
	if httpErr != nil {
		err <- httpErr
	} else {
		res <- response.StatusCode
	}
}
