-- +migrate Down

DROP TABLE IF EXISTS projects CASCADE;
DROP FUNCTION IF EXISTS set_updated_at();
