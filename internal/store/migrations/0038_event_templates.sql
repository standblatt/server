-- A club's event templates are events of their own that take no registrations; new events are copied from them.
ALTER TABLE events ADD COLUMN template boolean NOT NULL DEFAULT false;
