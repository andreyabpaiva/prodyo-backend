-- +migrate Down

-- Drop index
DROP INDEX IF EXISTS idx_actions_status;

-- Drop constraint
ALTER TABLE actions
DROP CONSTRAINT IF EXISTS actions_status_check;

-- Drop status column
ALTER TABLE actions
DROP COLUMN IF EXISTS status;
