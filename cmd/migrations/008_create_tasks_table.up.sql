-- +migrate Up

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    iteration_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    assignee_id UUID,
    status VARCHAR(50) NOT NULL DEFAULT 'NotStarted',
    timer TIMESTAMPTZ,
    parent_task_id UUID, -- For sub-tasks
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (iteration_id) REFERENCES iterations(id) ON DELETE CASCADE,
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    CHECK (status IN ('NotStarted', 'InProgress', 'Completed'))
);

CREATE INDEX IF NOT EXISTS idx_tasks_iteration_id ON tasks (iteration_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks (assignee_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_parent_task_id ON tasks (parent_task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks (created_at);

-- Create trigger for updated_at
DROP TRIGGER IF EXISTS trg_tasks_set_updated_at ON tasks;
CREATE TRIGGER trg_tasks_set_updated_at
BEFORE UPDATE ON tasks
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

