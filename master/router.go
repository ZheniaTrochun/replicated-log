package master

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Request struct {
	Message string `json:"message"`
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

	slog.Info(`Storing`, "message", request.Message)

	id, err := storeMessage(request.Message)

	if err != nil {
		slog.Error(`Error while storing message.`, "message", request.Message, "error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Stored successfully", "message", request.Message, "id", id)
	c.Status(200)
}
