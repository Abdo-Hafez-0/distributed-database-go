package handlers

import (
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}

	var req types.InsertRequest
	if err := parseBody(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" {
		WriteJSON(w, http.StatusBadRequest, false, "database is required")
		return
	}

	if req.Table == "" {
		WriteJSON(w, http.StatusBadRequest, false, "table is required")
		return
	}

	if len(req.Data) == 0 {
		WriteJSON(w, http.StatusBadRequest, false, "data is required")
		return
	}

	if database.DB == nil {
		WriteJSON(w, http.StatusInternalServerError, false, "database connection is not initialized")
		return
	}

	if _, err := database.InsertRecord(database.DB, req.Table, req.Data); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	repReq := types.ReplicationRequest{
		Operation: "insert",
		Database:  req.Database,
		Table:     req.Table,
		Data:      req.Data,
	}

	if err := replication.Replicate(repReq); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, "record inserted locally but replication failed: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, true, "record inserted successfully")
}
