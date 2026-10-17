CREATE TYPE helper_role AS ENUM ('scorer', 'ammunition', 'results');

CREATE TABLE helper_codes (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id    bigint NOT NULL,
    event_id     bigint NOT NULL,
    name         text NOT NULL,
    role         helper_role NOT NULL,
    range_id     bigint,
    lanes        integer[] NOT NULL DEFAULT '{}',
    token_hash   text NOT NULL UNIQUE,
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_used_at timestamptz,
    FOREIGN KEY (event_id, tenant_id) REFERENCES events(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE,
    CHECK ((role = 'scorer') = (range_id IS NOT NULL AND cardinality(lanes) > 0))
);

ALTER TABLE registrations ADD COLUMN ammo_issued_at timestamptz;
