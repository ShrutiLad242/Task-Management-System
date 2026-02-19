# Task Management System (Golang)

## Overview

This project is a scalable RESTful Task Management Service built using Go. It demonstrates clean architecture, JWT authentication, PostgreSQL persistence, background workers using goroutines and channels, and proper separation of concerns.

The system allows users to create, view, and delete tasks. Tasks are automatically processed and marked as completed by a background worker after a configurable delay.

---

## Features

### Core Features

- Create task
- Get all tasks
- Get task by ID
- Delete task
- PostgreSQL persistence
- Clean architecture (Controller → Service → Repository → DB)
- Background worker with goroutines and channels
- Automatic task completion after configurable delay
- JWT authentication
- Authorization (users can access only their own tasks)
- Environment-based configuration
- Graceful shutdown

---

## Task Lifecycle

Each task follows this lifecycle:

1. **Pending** - Task is created in pending state
2. **In Progress** - Task moves to in_progress (can be updated by service logic)
3. **Completed** - Task is automatically marked as completed after AUTO_COMPLETE_DELAY

---

## Prerequisites

Before running this project, ensure you have:

- **Go** 1.26.0 or higher - [Download](https://golang.org/dl/)
- **PostgreSQL** 15 or **Docker & Docker Compose** - [Download Docker](https://www.docker.com/products/docker-desktop)
- **Git** - [Download](https://git-scm.com/)
- A code editor (VS Code recommended)

---

## Installation & Setup

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd Task-Management-System
```

### Step 2: Set Up Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskdb

JWT_SECRET=your_super_secret_jwt_key_here_change_in_production

AUTO_COMPLETE_DELAY=1m
```

**Important**: 
- Change `JWT_SECRET` to a strong random value in production
- Ensure `DB_HOST=localhost` matches your PostgreSQL setup

### Step 3: Start PostgreSQL

#### Option A: Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

This will start a PostgreSQL container automatically. Verify it's running:

```bash
docker-compose ps
```

#### Option B: Using Local PostgreSQL

If you have PostgreSQL installed locally, ensure it's running and create the database:

```bash
psql -U postgres -c "CREATE DATABASE taskdb;"
```

### Step 4: Install Go Dependencies

```bash
go mod download
```

Verify dependencies are installed:

```bash
go mod verify
```

---

## Running the Application

### Option 1: Run with `go run` (Development)

```bash
go run cmd/main.go
```

You should see output like:

```
Task worker started
Server running on port 8080
```

### Option 2: Build and Run Binary

Build the application:

```bash
go build -o task-management cmd/main.go
```

Run the binary:

**On Windows:**
```bash
.\task-management.exe
```

**On macOS/Linux:**
```bash
./task-management
```

## Health Check

Once the application is running, verify it's working:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

## Troubleshooting

### Issue: Database Connection Failed

**Error:** `failed to connect to database`

**Solution:**
- Verify PostgreSQL is running: `docker-compose ps`
- Check `.env` has correct `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- Ensure database exists: `docker-compose exec postgres psql -U postgres -c "\l"`

### Issue: Port Already in Use

**Error:** `listen tcp :8080: bind: An attempt was made to use a socket in a way forbidden by its access policy`

**Solution:**
- Change `PORT` in `.env` to an available port (e.g., 8081)
- Or kill the process using the port

### Issue: Migration Errors

**Error:** `migration failed`

**Solution:**
- Check PostgreSQL is running
- Verify migrations folder exists and has SQL files
- Check database permissions

### Issue: JWT Token Errors

**Error:** `token is invalid`

**Solution:**
- Verify `JWT_SECRET` in `.env` matches across requests
- Ensure token is valid (24-hour expiration)

---

## Quick Test Workflow

### 1. Start the Application

```bash
go run cmd/main.go
```

### 2. Create a User (Login)

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response: `{"token":"eyJ..."}`

### 3. Create a Task

Replace `<TOKEN>` with the token from step 2:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Task",
    "description": "This is a test task"
  }'
```

### 4. Get All Tasks

```bash
curl -X GET http://localhost:8080/tasks \
  -H "Authorization: Bearer <TOKEN>"
```

### 5. Get Specific Task

Replace `<TASK_ID>` with the ID from step 3:

```bash
curl -X GET http://localhost:8080/tasks/<TASK_ID> \
  -H "Authorization: Bearer <TOKEN>"
```

### 6. Delete a Task

```bash
curl -X DELETE http://localhost:8080/tasks/<TASK_ID> \
  -H "Authorization: Bearer <TOKEN>"
```

---

## Development

### Project Structure

```
.
├── cmd/main.go                          # Application entry point
├── internal/
│   ├── auth/jwt.go                      # JWT authentication
│   ├── config/config.go                 # Configuration management
│   ├── controller/                      # HTTP handlers
│   ├── db/postgres.go                   # Database connection
│   ├── middleware/auth_middleware.go    # Auth middleware
│   ├── model/                           # Data models
│   ├── repository/                      # Data access layer
│   ├── service/                         # Business logic layer
│   └── worker/task-worker.go            # Background worker
├── migrations/                          # Database migrations
├── docker-compose.yml                   # Docker configuration
├── .env.example                         # Environment template
└── README.md
```

### Architecture

The project follows **Clean Architecture** principles:

```
HTTP Request
    ↓
Controller (HTTP handlers)
    ↓
Service (Business logic)
    ↓
Repository (Data access)
    ↓
Database (PostgreSQL)
```

### Running Tests

```bash
go test ./...
```

### Adding a New Feature

1. Define the model in `internal/model/`
2. Create repository methods in `internal/repository/`
3. Implement service logic in `internal/service/`
4. Add controller handlers in `internal/controller/`
5. Register routes in `cmd/main.go`

---

## Common Errors and Solutions

| Error | Cause | Solution |
|-------|-------|----------|
| `invalid user id: user does not exist` | Foreign key violation | Ensure user exists in database before creating tasks |
| `unauthorized` | User trying to access other's task | Only admins can see all tasks; users see only their own |
| `task not found` | Task ID doesn't exist | Use valid task ID from GET /tasks |
| `database connection refused` | PostgreSQL not running | Start Docker: `docker-compose up -d` |
| `listen tcp :8080: bind` | Port already in use | Change PORT in .env or kill process on port 8080 |

---

## Stopping the Application

### Stop the Server

Press `Ctrl+C` in the terminal

### Stop PostgreSQL (Docker)

```bash
docker-compose down
```

To remove volumes (delete data):

```bash
docker-compose down -v
```

---

## Environment Variables Reference

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `PORT` | 8080 | No | Server port |
| `DB_HOST` | localhost | No | PostgreSQL hostname |
| `DB_PORT` | 5432 | No | PostgreSQL port |
| `DB_USER` | postgres | No | PostgreSQL username |
| `DB_PASSWORD` | postgres | No | PostgreSQL password |
| `DB_NAME` | taskdb | No | Database name |
| `JWT_SECRET` | secret | Yes | Secret key for JWT signing (change in prod!) |
| `AUTO_COMPLETE_DELAY` | 1m | No | Auto-complete delay (e.g., 30s, 2m, 1h) |
