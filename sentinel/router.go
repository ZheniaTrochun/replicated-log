package sentinel

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
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
		log.Errorf("Error while decoding body: %s", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//log.Infof("Replicating message %s", string(bodyBytes))

	isDuplicate := sync(request.Id, request.Value, request.Timestamp)

	if isDuplicate {
		log.Warnf("Duplicate message replication: %s", request)
	} else {
		log.Infof("Replicated message %s successfully", request)
	}

	c.Status(201)
}

func (r Request) String() string {
	res, _ := json.Marshal(r)

	//return fmt.Sprintf(`{"id": "%d", "message": "%s", "timestamp": "%d"}`, r.Id, r.Value, r.Timestamp)
	return string(res)
}
