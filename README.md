# Distributed Database System Using Go

## Project Overview

This project is a distributed database system implemented using the Go programming language. The system is designed using a master-slave architecture where multiple nodes communicate over the network using HTTP APIs.

The project demonstrates core distributed systems concepts such as:
- Data replication
- Distributed communication
- Fault tolerance
- Dynamic database and table creation
- CRUD operations across distributed nodes

The system consists of:
- One Master Node
- Two Slave Nodes
- Replication mechanism
- Basic fault tolerance support
- MySQL databases for each node

---

# Objectives

The main objectives of the project are:

- Build a distributed database system using Go
- Allow nodes to communicate over the network
- Implement CRUD operations
- Replicate data between nodes
- Maintain system functionality even if a node fails
- Understand distributed systems architecture and communication

---

# System Architecture

```text
                +----------------------+
                |     Master Node      |
                |----------------------|
                | Write Operations     |
                | Create Database      |
                | Create Tables        |
                | Replication Manager  |
                +----------+-----------+
                           |
          ----------------------------------------
          |                                      |
+----------------------+          +----------------------+
|     Slave Node 1     |          |     Slave Node 2     |
|----------------------|          |----------------------|
| Replicated Database  |          | Replicated Database  |
| Read Queries         |          | Read Queries         |
+----------------------+          +----------------------+
```

---

# Technologies Used

## Backend
- Go (Golang)
- HTTP REST APIs
- MySQL

## Database
- MySQL Server

## Frontend
- HTML/CSS/JavaScript or React

## Tools
- Git & GitHub
- Postman
- VS Code

---

# Project Features

## Master Node Features
- Create database dynamically
- Create tables dynamically
- Insert records
- Update records
- Delete records
- Drop databases
- Replicate operations to slave nodes

## Slave Node Features
- Receive replication requests
- Store replicated data
- Execute read operations
- Support search queries

## Distributed Features
- Multi-node communication
- Data replication
- Basic fault tolerance
- Health checking
- Retry mechanism for failed replication

---

# API Endpoints

## Master APIs

### Create Database
POST /create-db

### Create Table
POST /create-table

### Insert Record
POST /insert

### Update Record
PUT /update

### Delete Record
DELETE /delete

### Drop Database
DELETE /drop-db

---

## Slave APIs

### Replication Endpoint
POST /replicate

### Select Records
GET /select

### Search Records
GET /search

---

# Project Structure

```text
distributed-database-go/
│
├── master/
│   ├── handlers/
│   ├── replication/
│   ├── database/
│   ├── models/
│   ├── routes/
│   └── main.go
│
├── slave1/
│   ├── handlers/
│   ├── database/
│   ├── models/
│   ├── routes/
│   └── main.go
│
├── slave2/
│   ├── handlers/
│   ├── database/
│   ├── models/
│   ├── routes/
│   └── main.go
│
├── shared/
│   ├── config/
│   ├── utils/
│   ├── constants/
│   └── types/
│
├── sql/
│   ├── schema/
│   ├── seed/
│   └── migrations/
│
├── docs/
│   ├── diagrams/
│   └── report/
│
├── gui/
│
├── data/
│
├── README.md
├── .gitignore
├── .env
└── go.mod
```

---

# Database Design

Each node uses its own MySQL database instance.

Example:
- Master Database → Port 3306
- Slave1 Database → Port 3307
- Slave2 Database → Port 3308

This design allows proper distributed replication between nodes.

---

# Replication Flow

```text
Client Request
      ↓
Master Node
      ↓
Master saves data locally
      ↓
Master sends replication request
      ↓
Slave1 applies operation
Slave2 applies operation
```

---

# Fault Tolerance

The system implements basic fault tolerance:
- If one slave node fails, the system continues running
- Failed replication requests are logged
- Retry mechanisms attempt replication again later

---

# Setup Instructions

## 1. Clone Repository

```bash
git clone <repository-url>
cd distributed-database-go
```

---

## 2. Install Dependencies

```bash
go mod tidy
```

---

## 3. Configure Environment Variables

Create `.env` file:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=distributed_db
```

---

## 4. Run Master Node

```bash
go run master/main.go
```

---

## 5. Run Slave Nodes

```bash
go run slave1/main.go
```

```bash
go run slave2/main.go
```

---

# Team Members Responsibilities

| Member | Responsibility |
|---|---|
| Member 1 | Master Node + APIs |
| Member 2 | Database Layer + MySQL Engine |
| Member 3 | Replication + Fault Tolerance |
| Member 4 | Slave Nodes + Cross-Technology Worker |
| Member 5 | GUI + Testing + Documentation |

---

# Future Improvements

Possible future enhancements:
- Leader election
- Automatic failover
- Authentication system
- Query optimization
- Docker deployment
- Load balancing
- Advanced distributed transactions

---

# Conclusion

This project demonstrates the implementation of a basic distributed database system using Go and MySQL. It introduces important distributed systems concepts such as node communication, replication, fault tolerance, and distributed data management.

The system is designed to be modular, scalable, and easy to extend with additional distributed features in the future.