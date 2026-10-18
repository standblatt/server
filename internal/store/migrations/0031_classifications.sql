-- Rankings over several programmes or attempts of an event; definition = domain.Classification.
CREATE TABLE classifications (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id  bigint NOT NULL,
    event_id   bigint NOT NULL,
    position   integer NOT NULL DEFAULT 0,
    definition jsonb NOT NULL,
    FOREIGN KEY (event_id, tenant_id) REFERENCES events(id, tenant_id) ON DELETE CASCADE
);
CREATE INDEX classifications_event_idx ON classifications (event_id, position);
CREATE TRIGGER classifications_event_open BEFORE INSERT OR UPDATE OR DELETE ON classifications FOR EACH ROW EXECUTE FUNCTION assert_event_open();
