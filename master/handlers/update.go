package handlers

import (
	"encoding/json"
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, false, "method not allowed")
		return
	}

	var req types.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" || req.Table == "" {
		writeJSON(w, http.StatusBadRequest, false, "database name and table name are required")
		return
	}

	if err := database.UpdateRecord(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	if err := replication.ReplicateUpdate(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, "record updated locally but replication failed")
		return
	}

	writeJSON(w, http.StatusOK, true, "record updated successfully")
}
