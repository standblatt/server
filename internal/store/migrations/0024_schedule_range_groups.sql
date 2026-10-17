-- The schedule of an event: single days with optional shooting times and an optional period (season) with
-- weekly shooting times. days stays as the expanded list of shooting days.
ALTER TABLE events ADD COLUMN schedule jsonb NOT NULL DEFAULT '{}';
ALTER TABLE events DISABLE TRIGGER events_archived;
UPDATE events SET schedule = jsonb_build_object('days', (SELECT jsonb_agg(jsonb_build_object('date', d::text) ORDER BY d) FROM unnest(days) d));
ALTER TABLE events ENABLE TRIGGER events_archived;

-- Lane groups of a range, e.g. 25 m on lanes 1-10 and 50 m on lanes 11-15 of one pistol range.
CREATE TABLE range_groups (
    id        bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id bigint NOT NULL,
    range_id  bigint NOT NULL,
    name      text NOT NULL,
    lanes     integer[] NOT NULL,
    position  integer NOT NULL,
    UNIQUE (range_id, name),
    UNIQUE (id, tenant_id),
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE
);

ALTER TABLE programmes ADD COLUMN range_group_id bigint,
    ADD FOREIGN KEY (range_group_id, tenant_id) REFERENCES range_groups(id, tenant_id);
