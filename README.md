# crong
Schedule and automate API calls.

# Setup

## Environment
```bash
AUTH_SECRET=your_secret_key_here
```

# Tools
## Cron
- https://github.com/go-co-op/gocron
- https://github.com/robfig/cron

## Execution
- https://github.com/hibiken/asynq

## Framework
- https://github.com/gofiber/fiber

## DB
- https://github.com/go-gorm/gorm
- https://github.com/sqlc-dev/sqlc
- https://github.com/golang-jwt/jwt

# Milestones

## 1. Execute Cron
1. Given a string cron schedule like "0 */3 * * *", I want to execute it programmatically.
2. Execute multiple of them in tandem.

## 2. Implement DB
1. Store the cron schedules in DB
2. Pull from DB and execute them.

## 3. Implement API to store new entries


# Structure
crong/
├── cmd/
│   ├── api/
│   │   └── main.go          # HTTP server
│   └── worker/
│       └── main.go          # cron runner
│
│   ├── scheduler/           # cron system
│   │   ├── scheduler.go
│   │   ├── runner.go
│   │   ├── loader.go        # loads jobs from DB
│   │   └── registry.go      # maps job types → functions
│
│   ├── workers/             # actual job implementations
│   │   ├── email_worker.go
│   │   ├── webhook_worker.go
│   │   └── cleanup_worker.go
│
└── go.mod

# Flow
worker startup
   ↓
load jobs from DB
   ↓
register jobs in cron engine
   ↓
trigger job → scheduler
   ↓
resolve handler from registry
   ↓
call service logic