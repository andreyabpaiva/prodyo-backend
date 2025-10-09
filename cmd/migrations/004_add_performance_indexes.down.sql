-- +migrate Down

-- Drop performance indexes
DROP INDEX IF EXISTS idx_projects_created_at_desc;
DROP INDEX IF EXISTS idx_projects_name_gin;
DROP INDEX IF EXISTS idx_projects_description_gin;

DROP INDEX IF EXISTS idx_users_created_at_desc;
DROP INDEX IF EXISTS idx_users_name_gin;
DROP INDEX IF EXISTS idx_users_email_gin;

DROP INDEX IF EXISTS idx_project_members_project_created;
DROP INDEX IF EXISTS idx_project_members_user_created;

DROP INDEX IF EXISTS idx_projects_created_name;
DROP INDEX IF EXISTS idx_users_created_name;
