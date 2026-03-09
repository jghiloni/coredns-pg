-- +goose Up
-- +goose StatementBegin
CREATE TYPE dns_record_type AS ENUM ('A','AAAA','CAA','CNAME','MX','NS','SOA','SRV','TXT');

CREATE TABLE zones
(
    fqdn TEXT NOT NULL UNIQUE PRIMARY KEY
);

CREATE TABLE records
(
    id UUID NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    name VARCHAR(63) NOT NULL,
    zone TEXT NOT NULL,
    ttl INTEGER NOT NULL DEFAULT 300,
    content JSONB NOT NULL DEFAULT '{}'::JSONB,
    record_type dns_record_type NOT NULL,
    CONSTRAINT ttl_is_positive CHECK (ttl > 0),
    FOREIGN KEY (zone) REFERENCES zones(fqdn)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE records;
DROP TABLE zones;
DROP TYPE dns_record_type;
-- +goose StatementEnd
