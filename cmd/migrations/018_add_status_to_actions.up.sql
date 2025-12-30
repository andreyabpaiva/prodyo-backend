-- +migrate Up

-- Add status column to actions table
ALTER TABLE actions
ADD COLUMN status VARCHAR(50) NOT NULL DEFAULT 'NotStarted';

-- Add check constraint to ensure valid status values
ALTER TABLE actions
ADD CONSTRAINT actions_status_check
CHECK (status IN ('NotStarted', 'InProgress', 'Completed'));

-- Create index on status for better query performance
CREATE INDEX IF NOT EXISTS idx_actions_status ON actions (status);
