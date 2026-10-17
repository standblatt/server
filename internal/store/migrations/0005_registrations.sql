CREATE TYPE registration_status AS ENUM ('registered', 'shooting', 'completed');

CREATE TABLE registrations (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id     bigint NOT NULL,
    event_id      bigint NOT NULL,
    shooter_id    bigint NOT NULL,
    club_id       bigint,
    age_category  age_category NOT NULL,
    lane          integer CHECK (lane > 0),
    status        registration_status NOT NULL,
    registered_at timestamptz NOT NULL DEFAULT now(),
    started_at    timestamptz,
    completed_at  timestamptz,
    FOREIGN KEY (event_id, tenant_id) REFERENCES events(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (shooter_id, tenant_id) REFERENCES shooters(id, tenant_id),
    FOREIGN KEY (club_id, tenant_id) REFERENCES clubs(id, tenant_id),
    UNIQUE (id, tenant_id)
);
CREATE UNIQUE INDEX registrations_lane_taken ON registrations (event_id, lane)
    WHERE lane IS NOT NULL AND status IN ('registered', 'shooting');
CREATE INDEX registrations_event_idx ON registrations (event_id, status);

CREATE TYPE entry_status AS ENUM ('pending', 'completed');

CREATE TABLE entries (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id       bigint NOT NULL,
    registration_id bigint NOT NULL,
    event_id        bigint NOT NULL,
    shooter_id      bigint NOT NULL,
    programme_id    bigint NOT NULL,
    status          entry_status NOT NULL DEFAULT 'pending',
    award_tier_id   bigint REFERENCES award_tiers(id),
    combined_series integer[],
    FOREIGN KEY (registration_id, tenant_id) REFERENCES registrations(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (programme_id, tenant_id) REFERENCES programmes(id, tenant_id),
    CONSTRAINT entries_programme_once UNIQUE (event_id, shooter_id, programme_id)
);
CREATE INDEX entries_registration_idx ON entries (registration_id);
CREATE INDEX entries_programme_idx ON entries (programme_id);
