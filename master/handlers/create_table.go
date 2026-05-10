package handlers

import (
	"encoding/json"
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, false, "method not allowed")
		return
	}

	var req types.CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" || req.Table == "" {
		writeJSON(w, http.StatusBadRequest, false, "database name and table name are required")
		return
	}

	if err := database.CreateTable(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	if err := replication.ReplicateCreateTable(req); err != nil {
		writeJSON(w, http.StatusInternalServerError, false, "table created locally but replication failed")
		return
	}

	writeJSON(w, http.StatusOK, true, "table created successfully")
}
