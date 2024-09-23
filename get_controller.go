package main

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"net/http"
)

type Service interface {
	GetAll() []string
}

type Controller struct {
	service Service
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Info("Retrieving list of stored messages")

	messages := c.service.GetAll()

	serializedRes, err := json.Marshal(messages)

	if err != nil {
		log.Errorf("Error serializing response: %s", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(200)
	_, err = w.Write(serializedRes)

	log.Infof("Stored messages: %s", string(serializedRes))

	if err != nil {
		log.Errorf("Error serializing response: %s", err)
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
}
