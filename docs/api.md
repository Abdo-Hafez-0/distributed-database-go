# Distributed Database API Documentation

## Base URLs

Master:
http://localhost:8080

Slave1:
http://localhost:8081

Slave2:
http://localhost:8082

---

# Master APIs

## Create Database

POST /create-db

Request:
```json
{
  "database": "shop"
}
```

Response:
```json
{
  "success": true,
  "message": "Database created successfully"
}
```

---

## Create Table

POST /create-table

Request:
```json
{
  "database": "shop",
  "table": "users",
  "columns": {
    "id": "INT PRIMARY KEY",
    "name": "VARCHAR(255)"
  }
}
```

---

## Insert Record

POST /insert

Request:
```json
{
  "database": "shop",
  "table": "users",
  "data": {
    "id": 1,
    "name": "Ali"
  }
}
```

---

## Update Record

PUT /update

---

## Delete Record

DELETE /delete

---

## Drop Database

DELETE /drop-db

---

# Slave APIs

## Replication Endpoint

POST /replicate

---

## Select Records

GET /select

---

## Search Records

GET /search