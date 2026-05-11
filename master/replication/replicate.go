package replication

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
	"fmt"
	"net/http"
	"sync"

	"distributed-database-go/shared/types"
)

// SlaveURLs holds the list of slave node endpoints
var SlaveURLs = []string{
	"http://localhost:8081/replicate",
	"http://localhost:8082/replicate",
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

// Replicate sends the replication request to all slave nodes
func Replicate(req types.ReplicationRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal replication request: %w", err)
	}

	for _, url := range SlaveURLs {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("failed to replicate to %s: %w", url, err)
		}
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("replication failed on %s with status %d", url, resp.StatusCode)
		}
	}

	return nil
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