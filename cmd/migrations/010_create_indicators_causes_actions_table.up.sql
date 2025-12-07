-- +migrate Up

CREATE TABLE IF NOT EXISTS indicators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    iteration_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (iteration_id) REFERENCES iterations(id) ON DELETE CASCADE,
    UNIQUE(iteration_id)
);

CREATE TABLE IF NOT EXISTS causes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indicator_id UUID NOT NULL,
    metric VARCHAR(50) NOT NULL,
    description TEXT,
    productivity_level VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE CASCADE,
    CHECK (metric IN ('WorkVelocity', 'ReworkIndex', 'InstabilityIndex')),
    CHECK (productivity_level IN ('Ok', 'Alert', 'Critical'))
);

CREATE TABLE IF NOT EXISTS actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    indicator_id UUID NOT NULL,
    cause_id UUID NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (indicator_id) REFERENCES indicators(id) ON DELETE CASCADE,
    FOREIGN KEY (cause_id) REFERENCES causes(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_indicators_iteration_id ON indicators (iteration_id);
CREATE INDEX IF NOT EXISTS idx_causes_indicator_id ON causes (indicator_id);
CREATE INDEX IF NOT EXISTS idx_causes_metric ON causes (metric);
CREATE INDEX IF NOT EXISTS idx_causes_productivity_level ON causes (productivity_level);
CREATE INDEX IF NOT EXISTS idx_actions_indicator_id ON actions (indicator_id);
CREATE INDEX IF NOT EXISTS idx_actions_cause_id ON actions (cause_id);

-- Create triggers for updated_at
DROP TRIGGER IF EXISTS trg_indicators_set_updated_at ON indicators;
CREATE TRIGGER trg_indicators_set_updated_at
BEFORE UPDATE ON indicators
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_causes_set_updated_at ON causes;
CREATE TRIGGER trg_causes_set_updated_at
BEFORE UPDATE ON causes
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_actions_set_updated_at ON actions;
CREATE TRIGGER trg_actions_set_updated_at
BEFORE UPDATE ON actions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

