package master

import (
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
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
		log.Errorf("Error while decoding body: %s", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Infof(`Storring message "%s"`, request.Message)

	_, err = storeMessage(request.Message)

	log.Infof(`Storred message "%s" successfully`, request.Message)

	if err != nil {
		log.Errorf(`Error while storring message "%s"`, err)
		c.JSON(500, gin.H{"error": err.Error()})
	} else {
		c.Status(200)
	}
}
