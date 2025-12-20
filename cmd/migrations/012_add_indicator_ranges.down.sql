-- +migrate Down

-- Drop trigger first
DROP TRIGGER IF EXISTS trg_indicator_ranges_set_updated_at ON indicator_ranges;

-- Drop indexes
DROP INDEX IF EXISTS idx_indicator_ranges_indicator_id;
DROP INDEX IF EXISTS idx_indicator_ranges_metric;

-- Remove calculated value columns from indicators
ALTER TABLE indicators
DROP COLUMN IF EXISTS velocity_value,
DROP COLUMN IF EXISTS rework_value,
DROP COLUMN IF EXISTS instability_value;

-- Drop the indicator_ranges table
DROP TABLE IF EXISTS indicator_ranges;

