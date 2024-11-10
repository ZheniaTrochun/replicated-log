package master

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Request struct {
	Message      string `json:"message"`
	WriteConcern int    `json:"writeConcern"`
}

func InitRouter(router *gin.Engine) {
	router.POST("/insert", insertHandler)
}

func insertHandler(c *gin.Context) {
	var request Request

	err := c.BindJSON(&request)

	if err != nil {
		slog.Error("Error while decoding body.", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	slog.Info(`Storing`, "message", request.Message, "write-concern", request.WriteConcern)

	id, err := storeMessage(request.Message, request.WriteConcern)

	if err != nil {
		slog.Error(`Error while storing message.`, "message", request.Message, "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Stored successfully", "id", id, "message", request.Message)
	c.Status(200)
}
