-- +migrate Up

-- Change timer column from TIMESTAMPTZ to BIGINT (stores seconds)
ALTER TABLE tasks ALTER COLUMN timer DROP DEFAULT;
ALTER TABLE tasks ALTER COLUMN timer TYPE BIGINT USING EXTRACT(EPOCH FROM timer)::BIGINT;
ALTER TABLE tasks ALTER COLUMN timer DROP NOT NULL;
COMMENT ON COLUMN tasks.timer IS 'Duration in seconds (e.g., 7200 for 2 hours)';
