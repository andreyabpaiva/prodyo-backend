-- +migrate Down

ALTER TABLE tasks DROP COLUMN IF EXISTS expected_time;
