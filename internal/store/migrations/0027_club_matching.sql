-- The club of the master data a booking request was matched to when it was confirmed. Requests confirmed before
-- are matched the way confirming did until now: the club's name regardless of case.
ALTER TABLE booking_requests ADD COLUMN club_id bigint,
    ADD FOREIGN KEY (club_id, tenant_id) REFERENCES clubs(id, tenant_id) ON DELETE SET NULL (club_id);
UPDATE booking_requests b SET club_id = (SELECT c.id FROM clubs c WHERE c.tenant_id = b.tenant_id AND lower(c.name) = lower(btrim(b.club))
    ORDER BY c.id LIMIT 1)
WHERE b.status = 'confirmed';

-- Merging two clubs moves registrations of closed events too: a transaction that sets standblatt.merge passes.
CREATE OR REPLACE FUNCTION assert_event_open() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    r  record;
    ev bigint;
BEGIN
    IF TG_OP = 'DELETE' THEN r := OLD; ELSE r := NEW; END IF;
    IF current_setting('standblatt.merge', true) = 'on' THEN
        RETURN r;
    END IF;
    CASE TG_TABLE_NAME
        WHEN 'results' THEN
            SELECT event_id INTO ev FROM entries WHERE id = r.entry_id;
        WHEN 'programme_series', 'programme_prices', 'award_tiers' THEN
            SELECT event_id INTO ev FROM programmes WHERE id = r.programme_id;
        WHEN 'award_limits' THEN
            SELECT p.event_id INTO ev FROM award_tiers t JOIN programmes p ON p.id = t.programme_id WHERE t.id = r.tier_id;
        ELSE
            ev := r.event_id;
    END CASE;
    IF EXISTS (SELECT 1 FROM events WHERE id = ev AND archived_at IS NOT NULL) THEN
        RAISE EXCEPTION 'event % is archived', ev USING ERRCODE = 'TLARC';
    END IF;
    RETURN r;
END $$;

CREATE OR REPLACE FUNCTION assert_event_unarchived() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF OLD.archived_at IS NOT NULL AND (TG_OP = 'DELETE' OR NEW.archived_at IS NOT NULL)
        AND coalesce(current_setting('standblatt.merge', true), '') <> 'on' THEN
        RAISE EXCEPTION 'event % is archived', OLD.id USING ERRCODE = 'TLARC';
    END IF;
    IF TG_OP = 'DELETE' THEN RETURN OLD; END IF;
    RETURN NEW;
END $$;
