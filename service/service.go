package service

import (
	"flag"
	"logger/helpers"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	addr    = flag.String("addr", ":8080", "The address to bind to")
	brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
)

type Config struct {
	SvcHost      string
	KafkaAddress string
	DbUser       string
	DbPassword   string
	DbHost       string
	DbName       string
}

type LoggerService struct {
}

func (s *LoggerService) Run(cfg Config) error {

	logResource := &LogResource{kafkaHelper: helpers.NewKafkaProducerHelper([]string{cfg.KafkaAddress})}
	r := gin.Default()
	r.Use(gin.Logger())

	r.POST("/log", logResource.CreateLog)

	r.Run()
	return nil
}
