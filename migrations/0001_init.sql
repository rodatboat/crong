-- +goose Up
CREATE TABLE folders (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    user_id INT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    verified BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Main job definition
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id INT NOT NULL,
    folder_id INT,

    method SMALLINT NOT NULL,
    headers JSONB,
    auth JSONB,
    cron TEXT NOT NULL,

    body TEXT,

    timezone TEXT NOT NULL DEFAULT 'UTC',

    timeout INT DEFAULT 30,
    enabled BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    -- FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE SET NULL
);

-- Minute(s) the job runs
CREATE TABLE schedule_minute (
    job_id INT NOT NULL,
    minute INT NOT NULL,
    PRIMARY KEY (job_id, minute),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Hour(s) the job runs
CREATE TABLE schedule_hour (
    job_id INT NOT NULL,
    hour INT NOT NULL,
    PRIMARY KEY (job_id, hour),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Day(s) of month the job runs
CREATE TABLE schedule_mday (
    job_id INT NOT NULL,
    mday INT NOT NULL,
    PRIMARY KEY (job_id, mday),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Day(s) of week the job runs
CREATE TABLE schedule_wday (
    job_id INT NOT NULL,
    wday INT NOT NULL,
    PRIMARY KEY (job_id, wday),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- Month(s) the job runs
CREATE TABLE schedule_month (
    job_id INT NOT NULL,
    month INT NOT NULL,
    PRIMARY KEY (job_id, month),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

CREATE TABLE job_executions (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,

    success BOOLEAN,
    status_code INT,
    duration_ms INT,
    url TEXT,
    batch_identifier TEXT,
    response_body TEXT,
    error TEXT,

    executed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS job_executions;
DROP TABLE IF EXISTS schedule_minute;
DROP TABLE IF EXISTS schedule_hour;
DROP TABLE IF EXISTS schedule_mday;
DROP TABLE IF EXISTS schedule_wday;
DROP TABLE IF EXISTS schedule_month;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS users;