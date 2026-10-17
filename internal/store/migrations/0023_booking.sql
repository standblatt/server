-- Online booking: external clubs ask for lanes in time slots through a public link per event; an admin
-- confirms the request, which turns its items into ordinary reservations and its shooters into master data.
ALTER TABLE events
    ADD COLUMN booking_code text UNIQUE,
    ADD COLUMN booking      jsonb NOT NULL DEFAULT '{}';

CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'declined', 'cancelled');

CREATE TABLE booking_requests (
    id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id      bigint NOT NULL,
    event_id       bigint NOT NULL,
    club           text NOT NULL,
    contact        text NOT NULL,
    email          text NOT NULL,
    phone          text NOT NULL DEFAULT '',
    note           text NOT NULL DEFAULT '',
    status         booking_status NOT NULL DEFAULT 'pending',
    decline_reason text NOT NULL DEFAULT '',
    manage_hash    text NOT NULL UNIQUE,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now(),
    decided_at     timestamptz,
    FOREIGN KEY (event_id, tenant_id) REFERENCES events(id, tenant_id) ON DELETE CASCADE,
    UNIQUE (id, tenant_id)
);
CREATE INDEX booking_requests_event_idx ON booking_requests (event_id);

CREATE TABLE booking_items (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id  bigint NOT NULL,
    request_id bigint NOT NULL,
    range_id   bigint NOT NULL,
    starts_at  timestamptz NOT NULL,
    ends_at    timestamptz NOT NULL CHECK (ends_at > starts_at),
    lanes      integer NOT NULL CHECK (lanes > 0),
    FOREIGN KEY (request_id, tenant_id) REFERENCES booking_requests(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE
);
CREATE INDEX booking_items_request_idx ON booking_items (request_id);

CREATE TABLE booking_shooters (
    id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id      bigint NOT NULL,
    request_id     bigint NOT NULL,
    firstname      text NOT NULL,
    lastname       text NOT NULL,
    birthday       date NOT NULL,
    license_number text NOT NULL DEFAULT '',
    shooter_id     bigint,
    FOREIGN KEY (request_id, tenant_id) REFERENCES booking_requests(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (shooter_id, tenant_id) REFERENCES shooters(id, tenant_id) ON DELETE SET NULL (shooter_id)
);
CREATE INDEX booking_shooters_request_idx ON booking_shooters (request_id);

ALTER TABLE reservations ADD COLUMN booking_item_id bigint REFERENCES booking_items(id) ON DELETE CASCADE;
