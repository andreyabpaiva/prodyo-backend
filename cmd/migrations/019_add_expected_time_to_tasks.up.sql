-- +migrate Up

ALTER TABLE tasks ADD COLUMN expected_time NUMERIC(10, 2) NOT NULL DEFAULT 0.0;

COMMENT ON COLUMN tasks.expected_time IS 'Expected time to complete the task in hours';
