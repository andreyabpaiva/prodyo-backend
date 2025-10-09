-- +migrate Down

-- Add email column back to projects table
ALTER TABLE projects ADD COLUMN email VARCHAR(255);

-- Recreate email index
CREATE INDEX IF NOT EXISTS idx_projects_email ON projects (email);

-- Drop project_members table
DROP TABLE IF EXISTS project_members CASCADE;
