ALTER TABLE events ADD COLUMN archived_at timestamptz;

-- Rows that belong to an archived event cannot be inserted, changed or deleted.
CREATE FUNCTION assert_event_open() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
    r  record;
    ev bigint;
BEGIN
    IF TG_OP = 'DELETE' THEN r := OLD; ELSE r := NEW; END IF;
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

CREATE TRIGGER registrations_event_open BEFORE INSERT OR UPDATE OR DELETE ON registrations FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER entries_event_open BEFORE INSERT OR UPDATE OR DELETE ON entries FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER results_event_open BEFORE INSERT OR UPDATE OR DELETE ON results FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER programmes_event_open BEFORE INSERT OR UPDATE OR DELETE ON programmes FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER programme_series_event_open BEFORE INSERT OR UPDATE OR DELETE ON programme_series FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER programme_prices_event_open BEFORE INSERT OR UPDATE OR DELETE ON programme_prices FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER award_tiers_event_open BEFORE INSERT OR UPDATE OR DELETE ON award_tiers FOR EACH ROW EXECUTE FUNCTION assert_event_open();
CREATE TRIGGER award_limits_event_open BEFORE INSERT OR UPDATE OR DELETE ON award_limits FOR EACH ROW EXECUTE FUNCTION assert_event_open();

-- An archived event can only be unarchived.
CREATE FUNCTION assert_event_unarchived() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF OLD.archived_at IS NOT NULL AND (TG_OP = 'DELETE' OR NEW.archived_at IS NOT NULL) THEN
        RAISE EXCEPTION 'event % is archived', OLD.id USING ERRCODE = 'TLARC';
    END IF;
    IF TG_OP = 'DELETE' THEN RETURN OLD; END IF;
    RETURN NEW;
END $$;

CREATE TRIGGER events_archived BEFORE UPDATE OR DELETE ON events FOR EACH ROW EXECUTE FUNCTION assert_event_unarchived();
