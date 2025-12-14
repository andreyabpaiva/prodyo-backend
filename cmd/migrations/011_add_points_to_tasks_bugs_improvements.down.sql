-- +migrate Down

-- Remove points column from improvements table
ALTER TABLE improvements DROP COLUMN IF EXISTS points;

-- Remove points column from bugs table
ALTER TABLE bugs DROP COLUMN IF EXISTS points;

-- Remove points column from tasks table
ALTER TABLE tasks DROP COLUMN IF EXISTS points;

