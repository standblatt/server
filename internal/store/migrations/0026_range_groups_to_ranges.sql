-- Lane groups become ranges of their own, numbered from 1: a range "Pistol" with 25 m on lanes 1-10 and 50 m on
-- lanes 11-15 turns into "Pistol 25 m" (the old range, so print agents stay, lanes 1-10) and "Pistol 50 m" (a new
-- range, lanes 1-5). Lanes outside every group move to a new range with the old name. Registrations, reservations,
-- scorer links (by their lowest lane), programmes, booking items, event ranges and print settings follow their lanes;
-- programmes tied to the range without a group lose the tie, so they still fit every lane they fitted before.
ALTER TABLE registrations DISABLE TRIGGER registrations_event_open;
ALTER TABLE programmes DISABLE TRIGGER programmes_event_open;

DO $$
DECLARE
    r      record;
    s      record;
    target bigint;
    old    bigint[];
    nm     text;
BEGIN
    FOR r IN SELECT * FROM ranges WHERE id IN (SELECT range_id FROM range_groups) ORDER BY id LOOP
        old := ARRAY(SELECT id FROM reservations WHERE range_id = r.id);
        UPDATE ranges SET name = '#' || id WHERE id = r.id;
        UPDATE programmes SET range_id = NULL WHERE range_id = r.id AND range_group_id IS NULL;
        -- The first group stays on the range; it comes last, when the other lanes have left it.
        FOR s IN
            SELECT g.id AS group_id, r.name || ' ' || g.name AS name, g.position = min(g.position) OVER () AS stays, g.position AS ord,
                ARRAY(SELECT DISTINCT l FROM unnest(g.lanes) l WHERE l BETWEEN 1 AND r.lane_count ORDER BY l) AS lanes
            FROM range_groups g WHERE g.range_id = r.id
            UNION ALL
            SELECT NULL, r.name, false, 2147483647, ARRAY(SELECT l FROM generate_series(1, r.lane_count) l
                WHERE NOT EXISTS (SELECT 1 FROM range_groups g WHERE g.range_id = r.id AND l = ANY(g.lanes)) ORDER BY l)
            ORDER BY stays, ord
        LOOP
            CONTINUE WHEN cardinality(s.lanes) = 0;
            nm := s.name;
            IF EXISTS (SELECT 1 FROM ranges WHERE tenant_id = r.tenant_id AND name = nm) THEN
                nm := nm || ' (' || coalesce(s.group_id, r.id) || ')';
            END IF;
            IF s.stays THEN
                target := r.id;
                UPDATE ranges SET name = nm, lane_count = cardinality(s.lanes),
                    disabled_lanes = ARRAY(SELECT array_position(s.lanes, d) FROM unnest(r.disabled_lanes) d WHERE d = ANY(s.lanes) ORDER BY 1)
                    WHERE id = r.id;
                -- Two steps, so the lane index never sees two registrations on one lane.
                UPDATE registrations SET lane = lane + 1000 WHERE range_id = r.id;
                UPDATE registrations SET lane = array_position(s.lanes, lane - 1000),
                    range_id = CASE WHEN lane - 1000 = ANY(s.lanes) THEN range_id END
                    WHERE range_id = r.id;
            ELSE
                INSERT INTO ranges (tenant_id, name, lane_count, disabled_lanes)
                    VALUES (r.tenant_id, nm, cardinality(s.lanes),
                        ARRAY(SELECT array_position(s.lanes, d) FROM unnest(r.disabled_lanes) d WHERE d = ANY(s.lanes) ORDER BY 1))
                    RETURNING id INTO target;
                UPDATE registrations SET range_id = target, lane = array_position(s.lanes, lane) WHERE range_id = r.id AND lane = ANY(s.lanes);
                INSERT INTO event_ranges (tenant_id, event_id, range_id, position)
                    SELECT tenant_id, event_id, target, position FROM event_ranges WHERE range_id = r.id;
                INSERT INTO print_settings (tenant_id, range_id, document, printer, options)
                    SELECT tenant_id, target, document, printer, options FROM print_settings WHERE range_id = r.id;
            END IF;
            IF s.group_id IS NOT NULL THEN
                UPDATE programmes SET range_id = target WHERE range_group_id = s.group_id;
            END IF;
            INSERT INTO reservations (tenant_id, range_id, lanes, starts_at, ends_at, lead_minutes, description, created_at, booking_item_id)
                SELECT tenant_id, target, ARRAY(SELECT array_position(s.lanes, l) FROM unnest(lanes) l WHERE l = ANY(s.lanes) ORDER BY 1),
                    starts_at, ends_at, lead_minutes, description, created_at, booking_item_id
                FROM reservations WHERE id = ANY(old) AND lanes && s.lanes;
            UPDATE booking_items b SET range_id = target WHERE b.range_id = r.id AND target <> r.id
                AND EXISTS (SELECT 1 FROM reservations x WHERE x.booking_item_id = b.id AND x.range_id = target);
            UPDATE helper_codes SET range_id = target,
                lanes = ARRAY(SELECT array_position(s.lanes, l) FROM unnest(lanes) l WHERE l = ANY(s.lanes) ORDER BY 1)
                WHERE range_id = r.id AND (SELECT min(l) FROM unnest(lanes) l) = ANY(s.lanes);
        END LOOP;
        DELETE FROM reservations WHERE id = ANY(old);
    END LOOP;
END $$;

UPDATE event_ranges er SET position = n.rn
FROM (SELECT event_id, range_id, row_number() OVER (PARTITION BY event_id ORDER BY position, range_id) AS rn FROM event_ranges) n
WHERE er.event_id = n.event_id AND er.range_id = n.range_id;

ALTER TABLE programmes DROP COLUMN range_group_id;
DROP TABLE range_groups;

ALTER TABLE registrations ENABLE TRIGGER registrations_event_open;
ALTER TABLE programmes ENABLE TRIGGER programmes_event_open;
