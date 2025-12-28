-- +migrate Down

DROP INDEX IF EXISTS idx_actions_end_at;
DROP INDEX IF EXISTS idx_actions_start_at;
DROP INDEX IF EXISTS idx_actions_assignee_id;

ALTER TABLE actions
DROP COLUMN assignee_id,
DROP COLUMN end_at,
DROP COLUMN start_at;
