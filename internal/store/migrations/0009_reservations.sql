CREATE TABLE reservations (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id    bigint NOT NULL,
    range_id     bigint NOT NULL,
    lanes        integer[] NOT NULL CHECK (cardinality(lanes) > 0),
    starts_at    timestamptz NOT NULL,
    ends_at      timestamptz NOT NULL CHECK (ends_at > starts_at),
    lead_minutes integer NOT NULL DEFAULT 0 CHECK (lead_minutes >= 0),
    description  text NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE
);
CREATE INDEX reservations_range_idx ON reservations (range_id, ends_at);
