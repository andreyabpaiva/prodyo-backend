-- +migrate Up

-- Add points column to tasks table
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 1;

-- Add points column to bugs table
ALTER TABLE bugs ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 1;

-- Add points column to improvements table
ALTER TABLE improvements ADD COLUMN IF NOT EXISTS points INTEGER NOT NULL DEFAULT 1;

-- Update all existing records to have points = 1 (in case there are any NULL values)
UPDATE tasks SET points = 1 WHERE points IS NULL;
UPDATE bugs SET points = 1 WHERE points IS NULL;
UPDATE improvements SET points = 1 WHERE points IS NULL;

