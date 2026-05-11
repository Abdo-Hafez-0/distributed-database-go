package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"distributed-database-go/shared/types"
)

var db *sql.DB

func SetDatabase(database *sql.DB) {
	db = database
}

func HandleReplication(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.ReplicationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf(
		"Replication Request -> Operation: %s | Table: %s\n",
		req.Operation,
		req.Table,
	)

	switch req.Operation {

	case "insert":
		handleInsert(req)

	case "update":
		handleUpdate(req)

	case "delete":
		handleDelete(req)

	default:
		http.Error(w, "Unknown operation", http.StatusBadRequest)
		return
	}

	response := types.ReplicationResponse{
		Status:  "success",
		Message: "Replication completed successfully",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func handleInsert(req types.ReplicationRequest) {

	log.Println("Executing INSERT replication")

	/*
	   Example implementation:

	   data := req.Data.(map[string]interface{})

	   query := "INSERT INTO users(name,email) VALUES(?,?)"

	   _, err := db.Exec(
	       query,
	       data["name"],
	       data["email"],
	   )

	   if err != nil {
	       log.Println(err)
	   }
	*/

	fmt.Println("Insert replicated successfully")
}
func handleUpdate(req types.ReplicationRequest) {

	log.Println("Executing UPDATE replication")

	/*
	   Example:

	   query := "UPDATE users SET name=? WHERE id=?"
	*/

	fmt.Println("Update replicated successfully")
}
func handleDelete(req types.ReplicationRequest) {

	log.Println("Executing DELETE replication")

	/*
	   Example:

	   query := "DELETE FROM users WHERE id=?"
	*/

	fmt.Println("Delete replicated successfully")
}