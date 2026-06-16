## Implementation Plan: DB-driven Cron HTTP Execution System (Fiber + Go)

### 1. Objective

Add a background scheduling system to an existing Go + Fiber application that:

* Stores cron jobs in a database
* Executes HTTP requests (GET/POST/PATCH/DELETE) on schedule
* Records execution history
* Runs non-blocking inside a single service process
* Enforces max request timeout of 30 seconds
* Uses bounded concurrency to prevent resource exhaustion

---

## 2. System Architecture

### Components (single binary)

1. **Fiber API Layer**

   * CRUD for jobs
   * Query job runs/history
   * Updates job definitions in DB

2. **Scheduler Goroutine (tick-based)**

   * Runs every 60 seconds
   * Queries DB for due jobs (`next_run_at <= now`)
   * Pushes jobs into execution queue

3. **Worker Pool**

   * Fixed number of goroutines
   * Consumes jobs from channel
   * Executes HTTP requests with timeout
   * Writes execution results + updates schedule

4. **Database Layer**

   * Jobs table (definition + schedule state)
   * Job runs table (audit log)

---

## 3. Database Schema (required)

### jobs

* id (uuid / bigint)
* url (string)
* method (string: GET/POST/PATCH/DELETE)
* headers (json)
* body (text/json)
* schedule (cron expression string)
* next_run_at (timestamp)
* last_run_at (timestamp)
* created_at (timestamp)
* updated_at (timestamp)

### job_runs

* id
* job_id
* status (success/fail)
* http_status_code
* response_body (text, truncated if needed)
* error_message
* started_at
* finished_at
* duration_ms

---

## 4. Runtime Design

### 4.1 Worker Pool

* Fixed-size channel buffer (queue)
* Fixed number of workers (5–10 for 512MB container)

```go
jobQueue := make(chan Job, 100)
```

Workers:

* continuously read from queue
* execute HTTP job
* persist result

---

### 4.2 Scheduler Loop (non-blocking)

Runs independently:

* Every 60 seconds:

  * Query DB for due jobs
  * Push into jobQueue (non-blocking send with backpressure handling)

Rules:

* Must NOT execute jobs directly
* Must NOT block on HTTP execution
* Must NOT spawn unlimited goroutines

---

### 4.3 Job Fetch Logic

Query:

* `SELECT * FROM jobs WHERE next_run_at <= NOW()`

Optional (recommended improvement):

* Use row locking if multiple instances are possible:

  * `FOR UPDATE SKIP LOCKED`

---

## 5. Execution Model

### 5.1 HTTP Execution Requirements

Each job execution must:

* Respect 30-second maximum timeout
* Use context cancellation
* Use `http.Client{Timeout: 30s}`
* Bind request to context

---

### 5.2 Execution Flow

For each job:

1. Create context with 30s timeout
2. Build HTTP request (method + headers + body)
3. Execute request
4. Capture:

   * status code
   * response body (truncate if needed)
   * error (if any)
5. Persist job_runs record
6. Compute next_run_at using cron parser
7. Update jobs table

---

## 6. Concurrency Rules

### 6.1 Scheduler

* Must be a single goroutine
* Must never execute HTTP calls

### 6.2 Worker Pool

* Fixed concurrency (5–10 workers recommended)
* Backpressure via buffered channel

### 6.3 Job Queue Behavior

* If queue is full:

  * drop job OR log and retry next cycle (configurable behavior)

---

## 7. Non-blocking Guarantee

System guarantees:

* Scheduler tick is independent of job execution time
* Job timeout does NOT block scheduler
* Max execution time is bounded by:

  ```
  (job_count / worker_count) × max_http_timeout
  ```

  not `job_count × timeout`

---

## 8. Scheduling Logic

### 8.1 Cron Parsing

* Use Go cron parser (`robfig/cron`)
* On job creation:

  * compute initial `next_run_at`

### 8.2 After execution:

* Recompute:

  ```
  next_run_at = cron.Next(now)
  ```

---

## 9. Fiber API Requirements

### Job endpoints

* `POST /jobs` → create job + compute next_run_at
* `GET /jobs` → list jobs
* `GET /jobs/:id` → job details
* `GET /jobs/:id/runs` → execution history
* `DELETE /jobs/:id` → remove job
* `PATCH /jobs/:id` → update job definition

---

## 10. Failure Handling

### HTTP failure cases

* timeout (30s)
* DNS failure
* non-2xx responses

All must:

* be recorded in job_runs
* NOT crash worker
* NOT stop scheduler

---

## 11. Resource Constraints (512MB / 1 vCPU)

Hard limits:

* Worker count: 5–10
* Queue buffer: 50–200
* HTTP timeout: 30s
* No unbounded goroutines

---

## 12. Optional Production Hardening (recommended later)

### 12.1 Multi-instance safety

* Add DB locking (`SKIP LOCKED`) OR Redis lock
* Prevent duplicate execution across containers

### 12.2 Retry strategy

* max retries: 2–3
* exponential backoff

### 12.3 Observability

* log job execution duration
* metrics:

  * queue depth
  * success/failure rate
  * execution latency

---

## 13. Implementation Order (for agent)

1. Add DB schema migrations
2. Implement job model + repository layer
3. Implement Fiber CRUD endpoints
4. Implement cron parsing + next_run_at calculation
5. Implement worker pool + job queue
6. Implement scheduler loop (tick every 60s)
7. Implement HTTP execution function (timeout + logging)
8. Implement job_runs persistence
9. Wire everything in `main()`
10. Add graceful shutdown (stop scheduler + workers)

---

## 14. Acceptance Criteria

System is correct if:

* Jobs execute at or near scheduled time
* API remains responsive under load
* No goroutine leaks under stress
* 100 concurrent jobs do not block scheduler
* Each job is capped at 30s execution
* Execution history is reliably stored
