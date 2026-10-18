-- A programme may be shot several times at an event (attempts); 0 means unlimited.
ALTER TABLE programmes
    ADD COLUMN attempts integer NOT NULL DEFAULT 1 CHECK (attempts >= 0),
    ADD COLUMN attempt_rule text NOT NULL DEFAULT 'best' CHECK (attempt_rule IN ('best', 'first', 'last')),
    ADD COLUMN award_counting_only boolean NOT NULL DEFAULT false;

-- Price of every attempt after the first; NULL = the same price.
ALTER TABLE programme_prices ADD COLUMN repeat_price numeric(8,2) CHECK (repeat_price >= 0);

ALTER TABLE entries ADD COLUMN attempt integer NOT NULL DEFAULT 1 CHECK (attempt > 0);
ALTER TABLE entries DROP CONSTRAINT entries_programme_once;
ALTER TABLE entries ADD CONSTRAINT entries_programme_once UNIQUE (event_id, shooter_id, programme_id, attempt);
