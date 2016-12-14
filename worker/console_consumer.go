package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/Shopify/sarama.v1"
	elastic "gopkg.in/olivere/elastic.v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type accessLogEntry struct {
	LogMessage []interface{}
}

var (
	brokerList = flag.String("brokers", "localhost:9092", "The comma separated list of brokers in the Kafka cluster")
	topic      = flag.String("topic", "access_log", "The topic to consume")
	partition  = flag.Int("partition", 0, "The partition to consume")
	offset     = flag.String("offset", "newest", "The offset to start with. Can be `oldest`, `newest`, or an actual offset")
	verbose    = flag.Bool("verbose", false, "Whether to turn on sarama logging")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {

	client, err := elastic.NewClient()

	if err != nil {
		// Handle error
	}

	flag.Parse()

	if *verbose {
		sarama.Logger = logger
	}

	var (
		initialOffset int64
		offsetError   error
	)
	switch *offset {
	case "oldest":
		initialOffset = sarama.OffsetOldest
	case "newest":
		initialOffset = sarama.OffsetNewest
	default:
		initialOffset, offsetError = strconv.ParseInt(*offset, 10, 64)
	}

	if offsetError != nil {
		logger.Fatalln("Invalid initial offset:", *offset)
	}

	c, err := sarama.NewConsumer(strings.Split(*brokerList, ","), nil)
	if err != nil {
		logger.Fatalln(err)
	}

	pc, err := c.ConsumePartition(*topic, int32(*partition), initialOffset)
	if err != nil {
		logger.Fatalln(err)
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		pc.AsyncClose()
	}()

	bulkRequest := client.Bulk()

	for msg := range pc.Messages() {
		fmt.Printf("Offset: %d\n", msg.Offset)
		fmt.Printf("Key:    %s\n", string(msg.Key))
		//fmt.Printf("Value:  %s\n", string(msg.Value))

		var accessLog accessLogEntry
		json.Unmarshal(msg.Value, &accessLog)
		//fmt.Println(accessLog)
		fmt.Println("working..")
		for i := range accessLog.LogMessage {

			indexReq := elastic.NewBulkIndexRequest().Index("events").Type("event").Id(strconv.FormatInt(msg.Offset, 10) + string(msg.Key) + "1").Doc(accessLog.LogMessage[i])
			bulkRequest = bulkRequest.Add(indexReq)

		}

		// messageUrl := "http://localhost:9200/events" + strconv.FormatInt(msg.Offset, 10) + string(msg.Key)
		// postToIndexer(messageUrl, msg.Value)

		bulkResponse, err := bulkRequest.Do(context.TODO())
		if err != nil {
			fmt.Println("erorr with error", err)

		}
		if bulkResponse == nil {
			fmt.Println("erorr with response", bulkResponse)

		}

	}

	if err := c.Close(); err != nil {
		fmt.Println("Failed to close consumer: ", err)
	}
}

func postToIndexer(url string, data []byte) {
	fmt.Println("URL:>", url)
	t := time.Now()
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("***** errror ********", err)
		return
	} else {
		defer resp.Body.Close()
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	elapsed := time.Since(t)
	fmt.Println("time %s", elapsed)

}
