package main

import (
	"github.com/apex/log"
	"net/http"
	"os"
	"replicated-log/master"
	"replicated-log/sentinel"
	"strings"
)

func main() {

	role := os.Getenv("ROLE")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	switch strings.ToLower(role) {
	case "master":
		sentinelsEnv := os.Getenv("SENTINELS")
		sentinelAddresses := strings.Split(sentinelsEnv, ",")

		masterService := master.NewLogMaster(sentinelAddresses)

		postController := master.NewController(masterService)
		getController := Controller{masterService}

		http.HandleFunc("POST /insert", postController.Insert)
		http.HandleFunc("GET /get-all", getController.GetAll)

	case "sentinel":
		sentinelService := sentinel.NewSentinel()

		postController := sentinel.NewController(sentinelService)
		getController := Controller{sentinelService}

		http.HandleFunc("POST "+sentinel.ReplicateEndpoint, postController.Replicate)
		http.HandleFunc("GET /get-all", getController.GetAll)
	}

	log.Infof("Starting server on port %s", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
