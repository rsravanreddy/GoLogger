package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger/api"
	"net/http"
	"sync"
	"time"
)

//  NodePath    string `json:"node_path"`
// 	ProcessPath string `json:"process_path"`
// 	TimeStamp   string `json:"time_stamp"`
// 	NodeId      int    `json:"node_id"`
// 	Action

func main() {
	var wg sync.WaitGroup
	wg.Add(10000)
	url := "http://localhost:8080/log"
	for i := 0; i < 1000; i++ {
		go testData(url, "/Application/chrome", "/usr/home/test.txt", "read", 1234, &wg)
	}
	// testData(url, "1", "message"+string(1), 1)
	wg.Wait()

}

func testData(url string, processPath string, nodePath string, action string, nodeId int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("URL:>", url)
	t := time.Now()
	sampleData := api.Log{NodePath: nodePath, ProcessPath: processPath, NodeId: nodeId, TimeStamp: t.Format("20060102150405"), Action: action}

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
