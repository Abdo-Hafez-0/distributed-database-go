package replication

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"distributed-database-go/shared/types"
)

// SlaveURLs holds the list of slave node endpoints
var SlaveURLs = []string{
	"http://localhost:8081/replicate",
	"http://localhost:8082/replicate",
}

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
