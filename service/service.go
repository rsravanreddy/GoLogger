package service

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"logger/helpers"
	"os"
	"strings"
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

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))
	kafkaHelper := &helpers.KafkaHelper{}
	logResource := &LogResource{AccessLogProducer: kafkaHelper.NewAccessLogProducer([]string{cfg.KafkaAddress})}
	r := gin.Default()
	r.Use(gin.Logger())

	r.POST("/log", logResource.CreateLog)

	lockScreenResouce := &LockScreenResouce{}

	r.GET("/lockScreen/", lockScreenResouce.LockScreen)

	r.Run()
	return nil
}

// func newAccessLogProducer(brokerList []string, cfg Config) sarama.AsyncProducer {s

// 	// For the access log, we are looking for AP semantics, with high throughput.
// 	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
// 	config := sarama.NewConfig()
// 	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
// 	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
// 	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

// 	producer, err := sarama.NewAsyncProducer([]string{"sravans-mbp-2:9092"}, config)
// 	if err != nil {
// 		log.Fatalln("Failed to start Sarama producer:", err)
// 	}

// 	// We will just log to STDOUT if we're not able to produce messages.
// 	// Note: messages will only be returned here after all retry attempts are exhausted.
// 	go func() {
// 		for err := range producer.Errors() {
// 			log.Println("Failed to write access log entry:", err)
// 		}
// 	}()

// 	return producer
// }
