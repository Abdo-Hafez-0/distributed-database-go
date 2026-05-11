package routes

import (
	"net/http"

	"distributed-database-go/master/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create-db", handlers.CreateDatabaseHandler)
	mux.HandleFunc("/create-table", handlers.CreateTableHandler)
	mux.HandleFunc("/insert", handlers.InsertHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete", handlers.DeleteHandler)
	mux.HandleFunc("/drop-db", handlers.DropDatabaseHandler)
}
