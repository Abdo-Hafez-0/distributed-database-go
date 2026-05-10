package handlers

import (
	"encoding/json"
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func DropDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, false, "method not allowed")
		return
	}

	var req types.DropDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" {
		writeJSON(w, http.StatusBadRequest, false, "database name is required")
		return
	}

	if err := database.DropDatabase(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	if err := replication.ReplicateDropDatabase(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, "database dropped locally but replication failed")
		return
	}

	writeJSON(w, http.StatusOK, true, "database dropped successfully")
}
