package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"os"
	"replicated-log/base"
	"replicated-log/master"
	"replicated-log/repository"
	"replicated-log/sentinel"
	"strings"
)

var roleEnvVariable = "ROLE"
var portEnvVariable = "PORT"
var grpcPortEnvVariable = "GRPC_PORT"
var sentinelsEnvVariable = "SENTINELS"

var roleMaster = "master"
var roleSentinel = "sentinel"

func main() {
	role := getRole()

	repository.InitDataStore()

	router := gin.New()

	switch strings.ToLower(role) {
	case roleMaster:
		slog.Info("Starting MASTER application...")

		sentinelAddresses := getSentinelAddresses()
		slog.Info("Retrieved", "sentinels", sentinelAddresses)

		master.InitLogMasterService(sentinelAddresses)
		master.InitRouter(router)
	case roleSentinel:
		slog.Info("Starting SENTINEL application...")

		server := grpc.NewServer()

		sentinel.InitServer(server)

		serverPort := getGrpcPort()
		lis, err := net.Listen("tcp", ":"+serverPort)
		if err != nil {
			slog.Error("Failed to start tcp listener.", "port", serverPort, "error", err)
			panic(err)
		}

		go func() {
			slog.Info("Starting gRPC server on port " + serverPort)
			err = server.Serve(lis)
			if err != nil {
				slog.Error("Failed to start sentinel gRPC server.", "error", err)
				panic(err)
			}
		}()
	}

	base.InitRouter(router)

	serveUrl := "0.0.0.0:" + getPort()

	slog.Info("Starting HTTP server on " + serveUrl)
	err := router.Run(serveUrl)

	if err != nil {
		slog.Error("Failed to start server.", "error", err)
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

func getGrpcPort() string {
	port := os.Getenv(grpcPortEnvVariable)

	if port == "" {
		port = "9090"
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
