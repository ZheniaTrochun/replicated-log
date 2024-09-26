package sentinel

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Request struct {
	Id        int    `json:"id"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

const replicateEndpoint = "/replicate-item"

func InitRouter(router *gin.Engine) {
	router.POST(replicateEndpoint, replicateHandler)
}

func replicateHandler(c *gin.Context) {
	var request Request
	err := c.BindJSON(&request)

	if err != nil {
		slog.Error("Error while decoding body.", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	isDuplicate := sync(request.Id, request.Value, request.Timestamp)

	if isDuplicate {
		slog.Warn("Duplicate item replication,", "message", request)
	} else {
		slog.Info("Replicated item successfully,", "message", request)
	}

	c.Status(201)
}

func (r Request) String() string {
	res, _ := json.Marshal(r)

	return string(res)
}
