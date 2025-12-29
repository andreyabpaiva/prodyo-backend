-- +migrate Up

-- Refactor causes and actions to use indicator_range_id instead of indicator_id
-- This changes them from iteration-specific to project-level

-- Delete all existing data from actions and causes tables
-- This is necessary because the old data references indicator IDs which are incompatible
-- with the new indicator_range_id structure
DELETE FROM actions;
DELETE FROM causes;

-- Drop the foreign key constraints and indexes
ALTER TABLE causes DROP CONSTRAINT IF EXISTS causes_indicator_id_fkey;
ALTER TABLE actions DROP CONSTRAINT IF EXISTS actions_indicator_id_fkey;
DROP INDEX IF EXISTS idx_causes_indicator_id;
DROP INDEX IF EXISTS idx_actions_indicator_id;

-- Rename the columns
ALTER TABLE causes RENAME COLUMN indicator_id TO indicator_range_id;
ALTER TABLE actions RENAME COLUMN indicator_id TO indicator_range_id;

-- Add new foreign key constraints linking to indicator_ranges table
ALTER TABLE causes
    ADD CONSTRAINT causes_indicator_range_id_fkey
    FOREIGN KEY (indicator_range_id) REFERENCES indicator_ranges(id) ON DELETE CASCADE;

ALTER TABLE actions
    ADD CONSTRAINT actions_indicator_range_id_fkey
    FOREIGN KEY (indicator_range_id) REFERENCES indicator_ranges(id) ON DELETE CASCADE;

-- Create new indexes
CREATE INDEX IF NOT EXISTS idx_causes_indicator_range_id ON causes (indicator_range_id);
CREATE INDEX IF NOT EXISTS idx_actions_indicator_range_id ON actions (indicator_range_id);
