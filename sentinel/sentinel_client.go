package sentinel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"replicated-log/repository"
	"strings"
)

type SentinelClient struct {
	address string
}

const ReplicateEndpoint = "/replicate-item"

func NewSentinelClient(address string) *SentinelClient {
	return &SentinelClient{address}
}

func (s *SentinelClient) SyncItem(item repository.Item, res chan int, err chan error) {
	request := Request{item.Id, item.Value, item.Timestamp}

	serializedRequest, serErr := json.Marshal(request)

	if serErr != nil {
		err <- serErr
	}

	bodyReader := bytes.NewReader([]byte(serializedRequest))

	requestUrl := withPrefix(s.address) + ReplicateEndpoint

	response, httpErr := http.Post(requestUrl, "application/json", bodyReader)
	if httpErr != nil {
		err <- httpErr
	} else {
		res <- response.StatusCode
	}
}

func withPrefix(address string) string {
	if strings.HasPrefix(address, "http://") ||
		strings.HasPrefix(address, "https://") {

		return address
	} else {
		return "http://" + address
	}
}
