CREATE TABLE results (
    entry_id      bigint NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
    series_number integer NOT NULL CHECK (series_number > 0),
    points        numeric(6,2) NOT NULL CHECK (points >= 0),
    mouches       integer CHECK (mouches >= 0),
    PRIMARY KEY (entry_id, series_number)
);
