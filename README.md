# crong

**Schedule and automate HTTP requests with precision cron expressions.**

Crong is a production-ready cron scheduling system that allows you to define, schedule, and execute HTTP API calls (GET, POST, PATCH, DELETE) on a precise schedule. It combines a REST API for job management with a background scheduler that runs jobs with bounded concurrency, comprehensive execution tracking, and detailed HTTP timing metrics.

## Features

- 🎯 **Cron-based scheduling** – Define jobs using standard cron expressions (`0 9 * * *`, etc.)
- 🔄 **HTTP execution** – Execute GET, POST, PATCH, DELETE requests on schedule with custom headers and auth
- 📊 **Execution history** – Track every job run with status codes, response bodies, and error messages
- ⏱️ **Timing metrics** – Capture DNS lookup, TCP connect, TLS handshake, and time-to-first-byte metrics
- 🛡️ **Bounded concurrency** – Fixed worker pool prevents resource exhaustion
- ⚙️ **Non-blocking architecture** – Scheduler and workers run independently without blocking the API
- 🔐 **User authentication** – JWT-based authentication for API endpoints
- 📁 **Job organization** – Group jobs into folders for easier management

## Tech Stack

- **Runtime:** Go 1.26+
- **API Framework:** [Fiber](https://gofiber.io/) (lightweight, fast HTTP framework)
- **Database:** PostgreSQL + [GORM](https://gorm.io/) ORM
- **Scheduling:** Custom tick-based scheduler with normalized schedule decomposition
- **Auth:** JWT tokens
- **Migrations:** [Goose](https://github.com/pressly/goose)

## Architecture

Crong runs as **two independent processes**:

### 1. API Server (`cmd/api/main.go`)
- REST endpoints for job CRUD operations
- User management and authentication
- Folder-based job organization
- Execution history queries

### 2. Scheduler (`cmd/scheduler/main.go`)
- Runs every 60 seconds (synchronized to minute boundary)
- Queries database for jobs due in the current minute
- Pushes jobs to a buffered queue
- Worker pool (default 5 workers) consumes and executes jobs
- Captures execution results and HTTP metrics

## Database Schema

Jobs are stored with **exploded schedule tables** for efficient per-minute queries:

```
jobs                 -- Job definitions
├── schedule_hours   -- Hour(s) the job runs
├── schedule_minutes -- Minute(s) the job runs
├── schedule_mdays   -- Day(s) of month
├── schedule_wdays   -- Day(s) of week
├── schedule_months  -- Month(s)
├── job_executions   -- Execution history
└── job_execution_stats -- HTTP timing metrics
```

## Setup

### Prerequisites

- Go 1.26+
- PostgreSQL 12+
- Docker (optional, for local DB)

### Environment

Create a `.env` file in the project root:

```bash
# Server
PORT=3000
AUTH_SECRET=your_secret_key_here

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=crong
DB_SSLMODE=disable
```

### Start PostgreSQL (Docker)

```bash
docker run --name postgres-crong \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=crong \
  -p 5432:5432 \
  -d postgres:latest
```

### Run Migrations

```bash
# DB name: postgres
# user: postgres
# password: password
# host: localhost:5435
# migrations folder: ./migrations
# schema: crong
goose -dir ./migrations postgres "postgres://postgres:password@localhost:5435/postgres?sslmode=disable&search_path=crong" up
```

## Running

### API Server

```bash
go run ./cmd/api/main.go
```

Server starts on `http://localhost:3000`

### Scheduler

In a separate terminal:

```bash
go run ./cmd/scheduler/main.go
```

The scheduler will:
1. Wait for the next minute boundary
2. Query the database every 60 seconds
3. Execute due jobs via the worker pool
4. Record execution results and metrics

## API Usage

### Authentication

All endpoints require a JWT token in the `Authorization` header:

```bash
Authorization: Bearer <token>
```

### Endpoints

#### Jobs

- `POST /api/jobs` – Create a new scheduled job
- `GET /api/jobs` – List all jobs
- `GET /api/jobs/:id` – Get job details
- `PUT /api/jobs/:id` – Update a job
- `DELETE /api/jobs/:id` – Delete a job

#### Example: Create a Job

```bash
POST /api/jobs
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Daily Report",
  "url": "https://api.example.com/report",
  "method": "POST",
  "schedule": "0 9 * * *",
  "timezone": "America/New_York",
  "headers": [
    { "key": "Content-Type", "value": "application/json" }
  ],
  "body": "{\"report_type\": \"daily\"}",
  "timeout": 30,
  "enabled": true
}
```

**Cron Expression Format:**
```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (0 = Sunday, 6 = Saturday)
│ │ │ │ │
│ │ │ │ │
* * * * *  -- Run every minute
0 9 * * *  -- Run at 9 AM every day
0 */3 * * * -- Run every 3 hours
0 0 * * 0   -- Run every Sunday at midnight
```

## Future Enhancements

- [ ] Retry logic with exponential backoff
- [ ] UI Application
- [ ] Search on execution history

## Development

### Project Structure

```
crong/
├── cmd/
│   ├── api/main.go         -- API server entry point
│   └── scheduler/main.go   -- Scheduler worker entry point
│   └── all/main.go         -- Both scheduler & API server entry point
├── internal/
│   ├── config/             -- Configuration loading
│   ├── database/           -- Database setup
│   ├── entities/           -- GORM entities (database models)
│   ├── models/             -- API request/response models
│   ├── services/           -- Business logic
│   ├── repositories/       -- Data access layer
│   ├── handlers/           -- HTTP request handlers
│   ├── middleware/         -- Auth, logging, etc.
│   ├── scheduler/          -- Cron scheduler logic
│   │   ├── scheduler.go    -- Main scheduler loop
│   │   ├── runner.go       -- Worker pool
│   │   └── loader.go       -- Job query logic
│   ├── routes/             -- Route definitions
│   └── utils/              -- Utilities
├── migrations/             -- Database migrations (Goose)
├── go.mod
└── README.md
```

### Running Tests

```bash
go test ./...
```

### Build

```bash
go build ./...
```

## License

GNU GPL v3.0