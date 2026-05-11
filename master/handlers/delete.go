package handlers

import (
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, "DELETE")
		return
	}

	var req types.DeleteRequest
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

	whereClause, whereArgs, err := buildWhereClause(req.Where)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, false, err.Error())
		return
	}

	if database.DB == nil {
		WriteJSON(w, http.StatusInternalServerError, false, "database connection is not initialized")
		return
	}

	if _, err := database.DeleteRecord(database.DB, req.Table, whereClause, whereArgs...); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	repReq := types.ReplicationRequest{
		Operation: "delete",
		Database:  req.Database,
		Table:     req.Table,
		Data: map[string]interface{}{
			"where": req.Where,
		},
	}

	if err := replication.Replicate(repReq); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, "record deleted locally but replication failed: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, true, "record deleted successfully")
}
