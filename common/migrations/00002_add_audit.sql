-- +goose Up
-- +goose StatementBegin
ALTER TABLE zones ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp;
ALTER TABLE zones ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp;
ALTER TABLE zones ADD COLUMN deleted_at TIMESTAMPTZ DEFAULT NULL;

ALTER TABLE records ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp;
ALTER TABLE records ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp;
ALTER TABLE records ADD COLUMN deleted_at TIMESTAMPTZ DEFAULT NULL;

CREATE INDEX IF NOT EXISTS idx_zones_created_at ON zones (created_at);
CREATE INDEX IF NOT EXISTS idx_zones_updated_at ON zones (updated_at);
CREATE INDEX IF NOT EXISTS idx_zones_deleted_at ON zones (deleted_at);

CREATE INDEX IF NOT EXISTS idx_records_created_at ON records (created_at);
CREATE INDEX IF NOT EXISTS idx_records_updated_at ON records (updated_at);
CREATE INDEX IF NOT EXISTS idx_records_deleted_at ON records (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_records_deleted_at;
DROP INDEX IF EXISTS idx_records_updated_at;
DROP INDEX IF EXISTS idx_records_created_at;

DROP INDEX IF EXISTS idx_zones_deleted_at;
DROP INDEX IF EXISTS idx_zones_updated_at;
DROP INDEX IF EXISTS idx_zones_created_at;

ALTER TABLE records DROP COLUMN deleted_at;
ALTER TABLE records DROP COLUMN updated_at;
ALTER TABLE records DROP COLUMN created_at;

ALTER TABLE zones DROP COLUMN deleted_at;
ALTER TABLE zones DROP COLUMN updated_at;
ALTER TABLE zones DROP COLUMN created_at;
-- +goose StatementEnd
