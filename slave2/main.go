package main

import (
	"log"
	"net/http"

	"distributed-database-go/slave2/handlers"
)

func main() {

	// Health Check Endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Replication Endpoint
	http.HandleFunc("/replicate", handlers.HandleReplication)

	log.Println("Slave2 running on port 8002")

	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		log.Fatal(err)
	}
}