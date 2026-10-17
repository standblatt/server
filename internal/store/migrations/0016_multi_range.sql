CREATE TABLE event_ranges (
    tenant_id bigint NOT NULL,
    event_id  bigint NOT NULL,
    range_id  bigint NOT NULL,
    position  integer NOT NULL,
    PRIMARY KEY (event_id, range_id),
    FOREIGN KEY (event_id, tenant_id) REFERENCES events(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id)
);
INSERT INTO event_ranges (tenant_id, event_id, range_id, position) SELECT tenant_id, id, range_id, 1 FROM events;

ALTER TABLE events ADD COLUMN days date[] NOT NULL DEFAULT '{}';
ALTER TABLE events DISABLE TRIGGER events_archived;
UPDATE events SET days = ARRAY[starts_on];
ALTER TABLE events ENABLE TRIGGER events_archived;
ALTER TABLE events ADD CHECK (cardinality(days) > 0), DROP COLUMN range_id;

ALTER TABLE registrations ADD COLUMN range_id bigint,
    ADD FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id);
ALTER TABLE registrations DISABLE TRIGGER registrations_event_open;
UPDATE registrations r SET range_id = er.range_id FROM event_ranges er WHERE er.event_id = r.event_id AND r.lane IS NOT NULL;
ALTER TABLE registrations ENABLE TRIGGER registrations_event_open;
ALTER TABLE registrations ADD CHECK ((lane IS NULL) = (range_id IS NULL));
DROP INDEX registrations_lane_taken;
CREATE UNIQUE INDEX registrations_lane_taken ON registrations (event_id, range_id, lane)
    WHERE lane IS NOT NULL AND status IN ('registered', 'shooting');

ALTER TABLE programmes ADD COLUMN range_id bigint,
    ADD FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id);
