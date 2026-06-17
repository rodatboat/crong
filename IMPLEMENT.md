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
   * Queries DB for due jobs
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

### Normalized Schedule Decomposition Schema

Jobs are stored with **exploded schedule tables** to enable efficient per-minute queries without cron parsing:

```
job                 -- Main job definition
├── job_hours       -- Hour(s) the job runs
├── job_minutes     -- Minute(s) the job runs
├── job_mdays       -- Day(s) of month the job runs
├── job_wdays       -- Day(s) of week the job runs
├── job_months      -- Month(s) the job runs
└── job_runs        -- Execution history
```

### jobs

* id (uuid / bigint)
* user_id (foreign key)
* url (string)
* method (string: GET/POST/PATCH/DELETE)
* headers (json)
* body (text/json)
* timezone (string, default: 'UTC')
* enabled (bool)
* created_at (timestamp)
* updated_at (timestamp)

**Key difference from cron string approach:** No schedule column. Schedule is decomposed into separate tables below.

### job_minutes

* job_id (composite PK)
* minute (integer 0-59, or -1 for "any")

Unique key: (job_id, minute)

### job_hours

* job_id (composite PK)
* hour (integer 0-23, or -1 for "any")

Unique key: (job_id, hour)

### job_mdays

* job_id (composite PK)
* mday (integer 1-31, or -1 for "any")

Unique key: (job_id, mday)

### job_wdays

* job_id (composite PK)
* wday (integer 0-6, or -1 for "any") [0=Sunday, 6=Saturday]

Unique key: (job_id, wday)

### job_months

* job_id (composite PK)
* month (integer 1-12, or -1 for "any")

Unique key: (job_id, month)

### job_runs

* id (primary key)
* job_id (foreign key)
* status (success/fail/timeout)
* http_status_code (integer, nullable)
* response_body (text, truncated if needed)
* error_message (text, nullable)
* started_at (timestamp)
* finished_at (timestamp)
* duration_ms (integer)

---

### Example: Job that runs every minute

Cron: `* * * * *` (every minute)

**jobs:**
```
id=123, user_id=1, url=..., enabled=1, timezone='UTC'
```

**job_minutes:** 60 rows (0-59, all minutes)
**job_hours:** 24 rows (0-23, all hours)
**job_mdays:** 1 row (-1, any day)
**job_wdays:** 1 row (-1, any day of week)
**job_months:** 1 row (-1, any month)

### Example: Job that runs at 9:00 AM every day

Cron: `0 9 * * *`

**jobs:**
```
id=456, user_id=1, url=..., enabled=1, timezone='UTC'
```

**job_minutes:** 1 row (minute=0)
**job_hours:** 1 row (hour=9)
**job_mdays:** 1 row (mday=-1)
**job_wdays:** 1 row (wday=-1)
**job_months:** 1 row (month=-1)

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

Runs independently every 60 seconds:

```go
for {
    select {
    case <-ticker.C:
        // Query for jobs matching current minute
        jobs := scheduler.QueryDueJobs(now)
        for _, job := range jobs {
            select {
            case jobQueue <- job:
                // Job queued
            default:
                // Queue full - log and skip (retry next minute)
                log.Warn("job queue full", "job_id", job.ID)
            }
        }
    case <-done:
        return
    }
}
```

**Query Logic:**

For the current minute (hour H, minute M, day D, month Mo, weekday W), execute:

```sql
SELECT j.* FROM jobs j
INNER JOIN job_hours jh ON j.id = jh.job_id AND (jh.hour = H OR jh.hour = -1)
INNER JOIN job_minutes jm ON j.id = jm.job_id AND (jm.minute = M OR jm.minute = -1)
INNER JOIN job_mdays jmd ON j.id = jmd.job_id AND (jmd.mday = D OR jmd.mday = -1)
INNER JOIN job_wdays jwd ON j.id = jwd.job_id AND (jwd.wday = W OR jwd.wday = -1)
INNER JOIN job_months jmo ON j.id = jmo.job_id AND (jmo.month = Mo OR jmo.month = -1)
WHERE j.enabled = 1
GROUP BY j.id
ORDER BY j.id ASC
```

Rules:

* Must NOT execute jobs directly
* Must NOT block on HTTP execution
* Must NOT spawn unlimited goroutines
* Must handle job queue backpressure (drop or log if full)

---

### 4.3 Job Fetch Logic (Stateless)

Key insight: **No state tracking needed**. Each minute's query inherently finds jobs due at that exact time.

* `GROUP BY job_id` deduplicates if a job matches multiple rules
* INNER JOINs guarantee all schedule conditions must match
* -1 values mean "match any" for that component
* Query is idempotent—same job re-picked if schedule allows

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
   * duration
5. Persist job_runs record
6. **No schedule update needed** — job will be re-picked by next minute's query if still due

**Key difference:** No cron parsing or next_run_at computation. The query-based approach is stateless and idempotent.

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

### 8.1 Schedule Representation

Instead of storing a cron string and parsing it:

* **Accept cron expression from API** (e.g., `0 9 * * *`)
* **Decompose into schedule tables** on job creation
* **Query via INNER JOINs** on each scheduler tick

Example: `0 9 * * *` (9 AM daily)

```
job_minutes: (job_id, 0)
job_hours: (job_id, 9)
job_mdays: (job_id, -1)  -- any day
job_wdays: (job_id, -1)  -- any day of week
job_months: (job_id, -1) -- any month
```

### 8.2 Cron Parsing Strategy

* On `POST /jobs` or `PATCH /jobs/:id`:
  * Parse cron expression using Go cron parser (`robfig/cron/v3`)
  * Extract all matching hours, minutes, mdays, wdays, months
  * Populate schedule tables (delete old entries, insert new ones)
  * Store -1 for "any" (e.g., day of month = -1 if not restricted)

### 8.3 Per-Minute Query

* No next_run_at timestamp needed
* Scheduler queries what's due **right now** based on current time
* Timezone-aware: convert current time to job's timezone before matching

---

## 9. Fiber API Requirements

### Job endpoints

* `POST /jobs` → create job + decompose cron into schedule tables
* `GET /jobs` → list jobs
* `GET /jobs/:id` → job details (include schedule decomposition)
* `GET /jobs/:id/runs` → execution history
* `DELETE /jobs/:id` → remove job + schedule table entries
* `PATCH /jobs/:id` → update job definition + re-decompose schedule if cron changed

### Request Body Format (POST/PATCH)

```json
{
  "url": "https://example.com/webhook",
  "method": "POST",
  "headers": {"Content-Type": "application/json"},
  "body": "{\"key\": \"value\"}",
  "schedule": "0 9 * * *",
  "timezone": "Europe/Berlin",
  "enabled": true
}
```

The `schedule` field is a cron expression. On job creation/update:
1. Parse the cron expression
2. Decompose into hours, minutes, mdays, wdays, months
3. Populate schedule tables with matching values (-1 for "any")

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

1. Add DB schema migrations (jobs + schedule tables + job_runs)
2. Implement job model + repository layer (includes schedule tables)
3. Implement cron decomposition function (parse cron → extract hours/minutes/mdays/wdays/months)
4. Implement Fiber CRUD endpoints (POST/GET/PATCH/DELETE jobs)
5. Implement schedule table population on job creation/update
6. Implement scheduler query (INNER JOINs across schedule tables for due jobs)
7. Implement worker pool + job queue
8. Implement HTTP execution function (timeout + logging)
9. Implement job_runs persistence
10. Wire everything in `main()` (start scheduler, start worker pool, start API)
11. Add graceful shutdown (stop scheduler + workers)

---

## 14. Acceptance Criteria

System is correct if:

* Jobs execute at or near scheduled time
* API remains responsive under load
* No goroutine leaks under stress
* 100 concurrent jobs do not block scheduler
* Each job is capped at 30s execution
* Execution history is reliably stored
