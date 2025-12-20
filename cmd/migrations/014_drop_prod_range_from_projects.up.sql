-- +migrate Up

-- Drop the prod_range column from projects table
-- Indicator ranges are now stored in the indicator_ranges table at project level
ALTER TABLE projects DROP COLUMN IF EXISTS prod_range;

