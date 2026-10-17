ALTER TABLE programmes
    ADD COLUMN prize_ranks        smallint NOT NULL DEFAULT 0 CHECK (prize_ranks BETWEEN 0 AND 10),
    ADD COLUMN prize_min_category smallint NOT NULL DEFAULT 0 CHECK (prize_min_category >= 0),
    ADD COLUMN prize_min_club     smallint NOT NULL DEFAULT 0 CHECK (prize_min_club >= 0);
