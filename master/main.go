package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func replicateToSlave() {
	data := map[string]interface{}{
		"operation": "insert",
		"database":  "shop",
		"table":     "users",
		"data": map[string]interface{}{
			"id":   1,
			"name": "Ali",
		},
	}

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(
		"http://localhost:8081/replicate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Println("Replication failed:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Replication sent successfully")
}

func main() {
	fmt.Println("Master running on port 8080")

	replicateToSlave()

	http.ListenAndServe(":8080", nil)
}