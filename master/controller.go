package master

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"io"
	"net/http"
)

type Request struct {
	Message string
}

type Controller struct {
	service *LogMaster
}

func NewController(service *LogMaster) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Insert(w http.ResponseWriter, r *http.Request) {
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

	log.Infof(`Storring message "%s"`, request.Message)

	_, err = c.service.StoreMessage(request.Message)

	log.Infof(`Storred message "%s" successfully`, request.Message)

	if err != nil {
		log.Errorf(`Error while storring message "%s"`, err)
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
}
