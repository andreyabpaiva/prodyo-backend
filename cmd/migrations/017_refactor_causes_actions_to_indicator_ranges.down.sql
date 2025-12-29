-- +migrate Down

-- Revert causes and actions back to using indicator_id instead of indicator_range_id

-- First, drop the foreign key constraints and indexes
ALTER TABLE causes DROP CONSTRAINT IF EXISTS causes_indicator_range_id_fkey;
ALTER TABLE actions DROP CONSTRAINT IF EXISTS actions_indicator_range_id_fkey;
DROP INDEX IF EXISTS idx_causes_indicator_range_id;
DROP INDEX IF EXISTS idx_actions_indicator_range_id;

-- Rename the columns back
ALTER TABLE causes RENAME COLUMN indicator_range_id TO indicator_id;
ALTER TABLE actions RENAME COLUMN indicator_range_id TO indicator_id;

-- Add back the original foreign key constraints linking to indicators table
ALTER TABLE causes
    ADD CONSTRAINT causes_indicator_id_fkey
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE CASCADE;

ALTER TABLE actions
    ADD CONSTRAINT actions_indicator_id_fkey
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE CASCADE;

-- Recreate original indexes
CREATE INDEX IF NOT EXISTS idx_causes_indicator_id ON causes (indicator_id);
CREATE INDEX IF NOT EXISTS idx_actions_indicator_id ON actions (indicator_id);
