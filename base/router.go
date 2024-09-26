package base

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	router.GET("/get-all", getAllHandler)
}

func getAllHandler(c *gin.Context) {
	log.Info("Retrieving list of stored messages")

	messages := getAll()

	log.Infof("Stored messages: %s", messages)

	c.JSON(200, messages)
}
