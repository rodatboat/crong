-- +goose Up

-- Composite index for listing jobs by schedule query
CREATE INDEX idx_jobs_enabled_timezone_id ON jobs(enabled, timezone, id);

-- Index on user_id for faster lookups and joins
CREATE INDEX idx_jobs_user_id ON jobs(user_id);

-- Index on folder_id for cascade operations and joins
CREATE INDEX idx_jobs_folder_id ON jobs(folder_id);

-- Indexes on schedule tables for faster joins
CREATE INDEX idx_schedule_minute_minute ON schedule_minute(minute);
CREATE INDEX idx_schedule_hour_hour ON schedule_hour(hour);
CREATE INDEX idx_schedule_mday_mday ON schedule_mday(mday);
CREATE INDEX idx_schedule_wday_wday ON schedule_wday(wday);
CREATE INDEX idx_schedule_month_month ON schedule_month(month);

-- +goose Down
DROP INDEX IF EXISTS idx_jobs_enabled_timezone_id;
DROP INDEX IF EXISTS idx_jobs_user_id;
DROP INDEX IF EXISTS idx_jobs_folder_id;
DROP INDEX IF EXISTS idx_schedule_minute_minute;
DROP INDEX IF EXISTS idx_schedule_hour_hour;
DROP INDEX IF EXISTS idx_schedule_mday_mday;
DROP INDEX IF EXISTS idx_schedule_wday_wday;
DROP INDEX IF EXISTS idx_schedule_month_month;
