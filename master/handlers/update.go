package handlers

import (
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		methodNotAllowed(w, "PUT")
		return
	}

	var req types.UpdateRequest
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

	whereClause, whereArgs, err := buildWhereClause(req.Where)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, false, err.Error())
		return
	}

	if database.DB == nil {
		WriteJSON(w, http.StatusInternalServerError, false, "database connection is not initialized")
		return
	}

	if _, err := database.UpdateRecord(database.DB, req.Table, req.Data, whereClause, whereArgs...); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	repData := map[string]interface{}{
		"where": req.Where,
		"data":  req.Data,
	}

	repReq := types.ReplicationRequest{
		Operation: "update",
		Database:  req.Database,
		Table:     req.Table,
		Data:      repData,
	}

	if err := replication.Replicate(repReq); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, "record updated locally but replication failed: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, true, "record updated successfully")
}
