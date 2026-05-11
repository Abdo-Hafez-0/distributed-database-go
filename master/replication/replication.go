package replication

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"distributed-database-go/shared/types"
)

var slaveNodes = []string{
	"http://localhost:8001",
	"http://localhost:8002",
}

const (
	maxRetries = 3
	retryDelay = 3 * time.Second
)

type FailedReplication struct {
	URL     string
	Request types.ReplicationRequest
	Retries int
}

var failedQueue []FailedReplication
var mutex sync.Mutex

func Replicate(req types.ReplicationRequest) {
	for _, slave := range slaveNodes {
		go replicateToSlave(slave, req)
	}
}

func replicateToSlave(slaveURL string, req types.ReplicationRequest) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Println("Replication marshal error:", err)
		return
	}

	endpoint := slaveURL + "/replicate"

	resp, err := http.Post(
		endpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Replication failed to %s\n", slaveURL)

		storeFailedReplication(slaveURL, req)

		return
	}

	log.Printf("Replication success -> %s\n", slaveURL)
}

func storeFailedReplication(url string, req types.ReplicationRequest) {
	mutex.Lock()
	defer mutex.Unlock()

	failedQueue = append(failedQueue, FailedReplication{
		URL:     url,
		Request: req,
		Retries: 0,
	})
}

func StartRetryWorker() {
	go func() {
		for {
			time.Sleep(retryDelay)

			mutex.Lock()

			var remaining []FailedReplication

			for _, item := range failedQueue {

				if item.Retries >= maxRetries {
					log.Printf("Max retries reached for %s\n", item.URL)
					continue
				}

				success := retryReplication(item)

				if !success {
					item.Retries++
					remaining = append(remaining, item)
				}
			}

			failedQueue = remaining

			mutex.Unlock()
		}
	}()
}

func retryReplication(item FailedReplication) bool {

	jsonData, _ := json.Marshal(item.Request)

	resp, err := http.Post(
		item.URL+"/replicate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil || resp.StatusCode != http.StatusOK {

		log.Printf(
			"Retry failed for %s attempt %d\n",
			item.URL,
			item.Retries+1,
		)

		return false
	}

	log.Printf("Retry success -> %s\n", item.URL)

	return true
}

func CheckNodeHealth(nodeURL string) bool {

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(nodeURL + "/health")

	if err != nil {
		log.Printf("Node unavailable: %s\n", nodeURL)
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}