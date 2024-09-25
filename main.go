package main

import (
	"github.com/apex/log"
	"net/http"
	"os"
	"replicated-log/base"
	"replicated-log/master"
	"replicated-log/repository"
	"replicated-log/sentinel"
	"strings"
)

func main() {
	role := os.Getenv("ROLE")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	repository.InitDataStore()

	switch strings.ToLower(role) {
	case "master":
		log.Info("Starting MASTER application...")

		sentinelAddresses := getSentinelAddresses()
		log.Infof("Sentinels: %s", sentinelAddresses)

		master.InitLogMasterService(sentinelAddresses)
		master.InitRouter()
	case "sentinel":
		log.Info("Starting SENTINEL application...")

		sentinel.InitRouter()
	}

	base.InitRouter()

	serveUrl := "0.0.0.0:" + port

	err := http.ListenAndServe(serveUrl, nil)

	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func getSentinelAddresses() []string {
	sentinelsEnv := os.Getenv("SENTINELS")

	if len(sentinelsEnv) == 0 {
		return make([]string, 0)
	} else {
		return strings.Split(sentinelsEnv, ",")
	}
}
