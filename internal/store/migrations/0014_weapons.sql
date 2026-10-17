CREATE TABLE weapons (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, id)
);
CREATE UNIQUE INDEX weapons_name_idx ON weapons (user_id, lower(name));

ALTER TABLE training_sessions ADD COLUMN weapon_id bigint,
    ADD FOREIGN KEY (user_id, weapon_id) REFERENCES weapons (user_id, id) ON DELETE SET NULL (weapon_id);
