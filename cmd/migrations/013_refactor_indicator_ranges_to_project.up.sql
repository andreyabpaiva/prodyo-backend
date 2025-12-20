-- +migrate Up

-- Drop the old indicator_ranges table and recreate with project_id
-- The ranges are now project-level, not indicator-level

-- First, drop triggers and indexes
DROP TRIGGER IF EXISTS trg_indicator_ranges_set_updated_at ON indicator_ranges;
DROP INDEX IF EXISTS idx_indicator_ranges_indicator_id;
DROP INDEX IF EXISTS idx_indicator_ranges_metric;

-- Drop the old table
DROP TABLE IF EXISTS indicator_ranges;

-- Create new indicator_ranges table linked to projects
-- Each project has its own set of ranges for each indicator type
CREATE TABLE IF NOT EXISTS indicator_ranges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL,
    indicator_type VARCHAR(50) NOT NULL,
    ok_min DECIMAL(10, 2) NOT NULL DEFAULT 0,
    ok_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    alert_min DECIMAL(10, 2) NOT NULL DEFAULT 0,
    alert_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    critical_min DECIMAL(10, 2) NOT NULL DEFAULT 0,
    critical_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CHECK (indicator_type IN ('SpeedPerIteration', 'ReworkPerIteration', 'InstabilityIndex')),
    UNIQUE(project_id, indicator_type)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_indicator_ranges_project_id ON indicator_ranges (project_id);
CREATE INDEX IF NOT EXISTS idx_indicator_ranges_indicator_type ON indicator_ranges (indicator_type);

-- Create trigger for updated_at on indicator_ranges
DROP TRIGGER IF EXISTS trg_indicator_ranges_set_updated_at ON indicator_ranges;
CREATE TRIGGER trg_indicator_ranges_set_updated_at
BEFORE UPDATE ON indicator_ranges
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

