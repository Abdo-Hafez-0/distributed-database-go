package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func replicateHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	fmt.Println("Replication request received:")
	fmt.Println(string(body))

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"success": true,
		"message": "Replication received",
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/replicate", replicateHandler)

	fmt.Println("Slave1 running on port 8081")

	http.ListenAndServe(":8081", nil)
}