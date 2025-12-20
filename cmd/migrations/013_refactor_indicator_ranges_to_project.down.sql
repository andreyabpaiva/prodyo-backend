-- +migrate Down

-- Drop trigger first
DROP TRIGGER IF EXISTS trg_indicator_ranges_set_updated_at ON indicator_ranges;

-- Drop indexes
DROP INDEX IF EXISTS idx_indicator_ranges_project_id;
DROP INDEX IF EXISTS idx_indicator_ranges_indicator_type;

-- Drop the project-based indicator_ranges table
DROP TABLE IF EXISTS indicator_ranges;

-- Recreate the old indicator-based indicator_ranges table
CREATE TABLE IF NOT EXISTS indicator_ranges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indicator_id UUID NOT NULL,
    metric VARCHAR(50) NOT NULL,
    ok_min DECIMAL(10, 2) NOT NULL DEFAULT 0,
    ok_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    alert_min DECIMAL(10, 2) NOT NULL DEFAULT 0,
    alert_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    critical_min DECIMAL(10, 2) NOT NULL DEFAULT 0,
    critical_max DECIMAL(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE CASCADE,
    CHECK (metric IN ('WorkVelocity', 'ReworkIndex', 'InstabilityIndex')),
    UNIQUE(indicator_id, metric)
);

CREATE INDEX IF NOT EXISTS idx_indicator_ranges_indicator_id ON indicator_ranges (indicator_id);
CREATE INDEX IF NOT EXISTS idx_indicator_ranges_metric ON indicator_ranges (metric);

CREATE TRIGGER trg_indicator_ranges_set_updated_at
BEFORE UPDATE ON indicator_ranges
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

