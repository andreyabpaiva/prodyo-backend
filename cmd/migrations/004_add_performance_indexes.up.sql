-- +migrate Up

-- Performance indexes for better query performance

-- Projects table indexes
CREATE INDEX IF NOT EXISTS idx_projects_created_at_desc ON projects (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_projects_name_gin ON projects USING gin (to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_projects_description_gin ON projects USING gin (to_tsvector('english', description));

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_created_at_desc ON users (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_name_gin ON users USING gin (to_tsvector('english', name));
CREATE INDEX IF NOT EXISTS idx_users_email_gin ON users USING gin (to_tsvector('english', email));

-- Project members junction table indexes
CREATE INDEX IF NOT EXISTS idx_project_members_project_created ON project_members (project_id, created_at);
CREATE INDEX IF NOT EXISTS idx_project_members_user_created ON project_members (user_id, created_at);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_projects_created_name ON projects (created_at DESC, name);
CREATE INDEX IF NOT EXISTS idx_users_created_name ON users (created_at DESC, name);

-- Partial indexes for active records (if you add soft deletes later)
-- CREATE INDEX IF NOT EXISTS idx_projects_active ON projects (created_at DESC) WHERE deleted_at IS NULL;
-- CREATE INDEX IF NOT EXISTS idx_users_active ON users (created_at DESC) WHERE deleted_at IS NULL;
