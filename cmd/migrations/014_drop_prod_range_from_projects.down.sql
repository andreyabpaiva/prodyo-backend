-- +migrate Down

-- Re-add the prod_range column to projects table
ALTER TABLE projects ADD COLUMN IF NOT EXISTS prod_range JSONB;

