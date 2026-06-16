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

CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id INT NOT NULL,
    folder_id INT,

    method SMALLINT NOT NULL,
    headers JSONB,
    auth JSONB,

    body TEXT,

    cron TEXT NOT NULL,
    timezone TEXT NOT NULL,

    timeout INT DEFAULT 30,
    enabled BOOLEAN DEFAULT TRUE,

    last_execution TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE job_executions (
    id SERIAL PRIMARY KEY,
    job_id INT NOT NULL,

    success BOOLEAN,
    status_code INT,
    duration_ms INT,
    response_body TEXT,
    error TEXT,

    executed_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS job_executions;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS folders;
DROP TABLE IF EXISTS users;