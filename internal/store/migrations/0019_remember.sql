ALTER TABLE sessions
    ADD COLUMN idle interval NOT NULL DEFAULT '7 days',
    ADD COLUMN ttl  interval NOT NULL DEFAULT '30 days';
