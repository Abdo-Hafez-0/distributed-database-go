package main

import (
	"log"
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/slave2/handlers"
)

func main() {
	db := database.MustConnect()
	handlers.SetDatabase(db)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/replicate", handlers.HandleReplication)
	http.HandleFunc("/select", handlers.HandleSelect)
	http.HandleFunc("/search", handlers.HandleSearch)

	log.Println("Slave2 running on port 8002")

	err := http.ListenAndServe(":8002", nil)
	if err != nil {
		log.Fatal(err)
	}
}
