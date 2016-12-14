package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger/api"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(10000)
	url := "http://localhost:8080/log"
	for i := 0; i < 1000; i++ {
		level := "error"
		if i%2 == 0 {
			level = "info"
		}
		go testData(url, level, "Russian guy vadim"+strconv.Itoa(i), i, &wg)
	}
	// testData(url, "1", "message"+string(1), 1)
	wg.Wait()

}

/*
LogLevel   string `json:"log_level"`
	LogPayload struct {
		Message   string `json:"message"`
		TimeStamp string `json:"time_stamp"`
		ErrorCode int    `json:"error_code"`
	} `json:"log_payload"`
}*/

func testData(url string, level string, msg string, er int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("URL:>", url)
	t := time.Now()
	sampleData := api.Log{LogLevel: level, LogPayload: api.LogPayloads{Message: msg, TimeStamp: t.Format("20060102150405"), ErrorCode: er}}
	b, err := json.Marshal(sampleData)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	elapsed := time.Since(t)
	fmt.Println("time %s", elapsed)

}
