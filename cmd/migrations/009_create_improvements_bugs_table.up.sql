-- +migrate Up

CREATE TABLE IF NOT EXISTS improvements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    assignee_id UUID,
    number INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(task_id, number)
);

CREATE TABLE IF NOT EXISTS bugs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id UUID NOT NULL,
    assignee_id UUID,
    number INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (assignee_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(task_id, number)
);

CREATE INDEX IF NOT EXISTS idx_improvements_task_id ON improvements (task_id);
CREATE INDEX IF NOT EXISTS idx_improvements_assignee_id ON improvements (assignee_id);
CREATE INDEX IF NOT EXISTS idx_improvements_number ON improvements (number);
CREATE INDEX IF NOT EXISTS idx_bugs_task_id ON bugs (task_id);
CREATE INDEX IF NOT EXISTS idx_bugs_assignee_id ON bugs (assignee_id);
CREATE INDEX IF NOT EXISTS idx_bugs_number ON bugs (number);

-- Create triggers for updated_at
DROP TRIGGER IF EXISTS trg_improvements_set_updated_at ON improvements;
CREATE TRIGGER trg_improvements_set_updated_at
BEFORE UPDATE ON improvements
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_bugs_set_updated_at ON bugs;
CREATE TRIGGER trg_bugs_set_updated_at
BEFORE UPDATE ON bugs
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

