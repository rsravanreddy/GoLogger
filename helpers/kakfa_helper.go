package helpers

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

type KafkaHelper struct {
	producer sarama.AsyncProducer
	consumer sarama.AsyncProducer
}

func NewKafkaProducerHelper(brokerList []string) *KafkaHelper {
	return &KafkaHelper{producer: NewAccessLogProducer(brokerList)}
}

func NewAccessLogProducer(brokerList []string) sarama.AsyncProducer {

	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}

func NewAccessLogConsumer(brokerList []string, topic string, partition int32, offSet int64) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	master, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
		return nil, err
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// How to decide partition, is it fixed value...?
	consumer, err := master.ConsumePartition(topic, partition, offSet)
	if err != nil {
		return nil, err
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	return consumer, nil
}

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

func (kh *KafkaHelper) ProduceAccesLogMessage(topic string, key string, entry []interface{}) {

	value := &accessLogEntry{
		LogMessage: entry,
	}

	kh.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: value,
	}

}
