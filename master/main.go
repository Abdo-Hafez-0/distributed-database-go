package main

import (
	"log"
	"net/http"

	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func main() {

	// Start Retry Worker
	replication.StartRetryWorker()

	// Health Check Endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("MASTER OK"))
	})

	// Temporary Testing Endpoint
	http.HandleFunc("/test-replication", func(w http.ResponseWriter, r *http.Request) {

		replication.Replicate(types.ReplicationRequest{
			Operation: "insert",
			Database:  "shop",
			Table:     "users",
			Data: map[string]interface{}{
				"id":   1,
				"name": "Abdo",
			},
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Replication sent successfully"))
	})

	log.Println("Master running on port 8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}