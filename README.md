# Internal Transfers API

This project implements an internal transfers application in Go, supporting account creation, balance querying, and atomic fund transfers between accounts, with persistence in PostgreSQL.

## Features

✅ Create accounts with initial balances  
✅ Query account balances  
✅ Submit transactions (fund transfers)  
✅ Consistent, atomic updates using PostgreSQL transactions  
✅ Dockerized environment with PostgreSQL  
✅ Schema migrations  
✅ Unit-tested services and handlers  
✅ Clean, maintainable code with clear separation of layers

## Running the Application
### Prerequisites
- Docker + Docker Compose
- Go 1.23+
### Steps
#### 1. Clone the repo:
```bash
git clone https://github.com/jasona122/internal-transfers.git
cd internal-transfers
```

#### 2. Set up environment variables
Create an .env file in the project root; use .env.example as reference

#### 3. Start Postgres:
```bash
docker compose up -d postgres
```
#### 4. Run migrations:
```bash
make migrate
```
#### 5. Build & run the app:
```bash
go build -o main ./cmd/main.go
./main
```
or in Docker:
```bash
docker compose up --build
```
## API Endpoints
[View in the Swagger Editor](https://editor.swagger.io/?url=https://raw.githubusercontent.com/jasona122/internal-transfers/docs/openapi.yml)

## Assumptions: 
- No authn/authz security
- Balance in user's account cannot be less than 0
- No need to encrypt user details in DB nor response
- Precision is limited to float64
- Database used is postgres
- Monetary values from client requests are in string format
  - Monetary values returned from server will be in float64 format, with its corresponding precision