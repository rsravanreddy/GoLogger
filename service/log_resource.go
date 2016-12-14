package service

import (
	"logger/api"
	"logger/helpers"

	"github.com/gin-gonic/gin"
)

type LogResource struct {
	kafkaHelper *helpers.KafkaHelper
}

func (lr *LogResource) CreateLog(c *gin.Context) {
	var accessLogs []interface{}
	if c.BindJSON(&accessLogs) != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	lr.kafkaHelper.ProduceAccesLogMessage("access_log", c.ClientIP(), accessLogs)

	c.JSON(200, api.NewSuccess("success", len(accessLogs)))
}
