package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gopkg.in/Shopify/sarama.v1"
	"logger/api"
)

type accessLogEntry struct {
	LogMessage []interface{} 
	encoded    []byte
	err        error
}

func (ale *accessLogEntry) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale)
	}
}

func (ale *accessLogEntry) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *accessLogEntry) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

type LogResource struct {
	AccessLogProducer sarama.AsyncProducer
}

func (lr *LogResource) CreateLog(c *gin.Context) {
	var accessLogs []interface{}
	if c.BindJSON(&accessLogs) != nil {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	entry := &accessLogEntry{
		LogMessage: accessLogs,
	}

	lr.AccessLogProducer.Input() <- &sarama.ProducerMessage{
		Topic: "access_log",
		Key:   sarama.StringEncoder(c.ClientIP()),
		Value: entry,
	}

	c.JSON(200, api.NewSuccess("success", len(accessLogs)))
}
