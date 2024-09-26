package base

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func InitRouter(router *gin.Engine) {
	router.GET("/get-all", getAllHandler)
}

func getAllHandler(c *gin.Context) {
	slog.Info("Retrieving list of stored messages")

	messages := getAll()

	slog.Info("Stored", "messages", messages)

	c.JSON(200, messages)
}
