-- +migrate Up

CREATE TABLE IF NOT EXISTS iterations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    number INTEGER NOT NULL,
    description TEXT,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE(project_id, number)
);

CREATE INDEX IF NOT EXISTS idx_iterations_project_id ON iterations (project_id);
CREATE INDEX IF NOT EXISTS idx_iterations_number ON iterations (number);
CREATE INDEX IF NOT EXISTS idx_iterations_start_at ON iterations (start_at);
CREATE INDEX IF NOT EXISTS idx_iterations_end_at ON iterations (end_at);

-- Create trigger for updated_at
DROP TRIGGER IF EXISTS trg_iterations_set_updated_at ON iterations;
CREATE TRIGGER trg_iterations_set_updated_at
BEFORE UPDATE ON iterations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

