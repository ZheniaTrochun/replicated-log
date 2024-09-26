package main

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"os"
	"replicated-log/base"
	"replicated-log/master"
	"replicated-log/repository"
	"replicated-log/sentinel"
	"strings"
)

var roleEnvVariable = "ROLE"
var portEnvVariable = "PORT"
var sentinelsEnvVariable = "SENTINELS"

var roleMaster = "master"
var roleSentinel = "sentinel"

func main() {
	role := getRole()

	repository.InitDataStore()

	router := gin.New()

	switch strings.ToLower(role) {
	case roleMaster:
		log.Info("Starting MASTER application...")

		sentinelAddresses := getSentinelAddresses()
		log.Infof("Sentinels: %s", sentinelAddresses)

		master.InitLogMasterService(sentinelAddresses)
		master.InitRouter(router)
	case roleSentinel:
		log.Info("Starting SENTINEL application...")

		sentinel.InitRouter(router)
	}

	base.InitRouter(router)

	serveUrl := "0.0.0.0:" + getPort()
	err := router.Run(serveUrl)

	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}

func getRole() string {
	role := os.Getenv(roleEnvVariable)

	if role == "" ||
		!(strings.ToLower(role) == roleMaster || strings.ToLower(role) == roleSentinel) {

		panic("Failed to read instance role")
	}

	return role
}

func getPort() string {
	port := os.Getenv(portEnvVariable)

	if port == "" {
		port = "8080"
	}

	return port
}

func getSentinelAddresses() []string {
	sentinelsEnv := os.Getenv(sentinelsEnvVariable)

	if len(sentinelsEnv) == 0 {
		return make([]string, 0)
	} else {
		return strings.Split(sentinelsEnv, ",")
	}
}
