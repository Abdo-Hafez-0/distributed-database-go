package handlers

import (
	"encoding/json"
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func CreateDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, false, "method not allowed")
		return
	}

	var req types.CreateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" {
		writeJSON(w, http.StatusBadRequest, false, "database name is required")
		return
	}

	if err := database.CreateDatabase(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	if err := replication.ReplicateCreateDatabase(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, "database created locally but replication failed")
		return
	}

	writeJSON(w, http.StatusOK, true, "database created successfully")
}
