package sentinel

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"io"
	"net/http"
)

type Request struct {
	Id        int
	Value     string
	Timestamp int64
}

func InitRouter() {
	endpointPattern := "POST " + ReplicateEndpoint

	http.HandleFunc(endpointPattern, replicateHandler)
}

func replicateHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.Errorf("Error while reading body: %s", err)
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	var request Request

	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		log.Errorf("Error while decoding body: %s", err)
		w.WriteHeader(400)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	//log.Infof("Replicating message %s", string(bodyBytes))

	isDuplicate := sync(request.Id, request.Value, request.Timestamp)

	if isDuplicate {
		log.Warnf("Duplicate message replication: %s", string(bodyBytes))
	} else {
		log.Infof("Replicated message %s successfully", string(bodyBytes))
	}

	w.WriteHeader(200)
}
