package types

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CreateDatabaseRequest struct {
	Database string `json:"database"`
}

type CreateTableRequest struct {
	Database string            `json:"database"`
	Table    string            `json:"table"`
	Columns  map[string]string `json:"columns"`
}
type UpdateRequest struct {
	Database string                 `json:"database"`
	Table    string                 `json:"table"`
	Where    map[string]interface{} `json:"where"`
	Data     map[string]interface{} `json:"data"`
}

type DeleteRequest struct {
	Database string                 `json:"database"`
	Table    string                 `json:"table"`
	Where    map[string]interface{} `json:"where"`
}

type DropDatabaseRequest struct {
	Database string `json:"database"`
}
type InsertRequest struct {
	Database string                 `json:"database"`
	Table    string                 `json:"table"`
	Data     map[string]interface{} `json:"data"`
}

type ReplicationRequest struct {
	Operation string                 `json:"operation"`
	Database  string                 `json:"database"`
	Table     string                 `json:"table"`
	Data      map[string]interface{} `json:"data"`
}
