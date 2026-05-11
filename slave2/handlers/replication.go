package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"distributed-database-go/shared/types"
)

func HandleReplication(w http.ResponseWriter, r *http.Request) {

	var req types.ReplicationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf(
		"Replication received: %s on table %s\n",
		req.Operation,
		req.Table,
	)

	// TODO:
	// execute actual DB operation here
	// insert/update/delete

	response := types.ReplicationResponse{
		Status:  "success",
		Message: "Replication completed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}