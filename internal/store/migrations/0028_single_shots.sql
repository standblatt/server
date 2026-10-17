-- Series can record every single shot; the series points are then the sum of the shots.
ALTER TABLE programme_series ADD COLUMN track_shots boolean NOT NULL DEFAULT false;
ALTER TABLE results ADD COLUMN shots numeric[];
