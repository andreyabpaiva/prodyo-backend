-- +migrate Up

-- Add password hash to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);

-- Create index for email lookups (for login)
CREATE INDEX IF NOT EXISTS idx_users_email_unique ON users (email) WHERE email IS NOT NULL;

