-- +migrate Down

-- Revert timer column from BIGINT back to TIMESTAMPTZ
ALTER TABLE tasks ALTER COLUMN timer TYPE TIMESTAMPTZ USING to_timestamp(timer);
COMMENT ON COLUMN tasks.timer IS NULL;
