-- +migrate Up

ALTER TABLE actions
ADD COLUMN start_at TIMESTAMPTZ,
ADD COLUMN end_at TIMESTAMPTZ,
ADD COLUMN assignee_id UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_actions_assignee_id ON actions (assignee_id);
CREATE INDEX IF NOT EXISTS idx_actions_start_at ON actions (start_at);
CREATE INDEX IF NOT EXISTS idx_actions_end_at ON actions (end_at);
