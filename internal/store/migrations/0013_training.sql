ALTER TABLE users
    ADD COLUMN email                  text UNIQUE CHECK (email = lower(email)),
    ADD COLUMN email_verified_at      timestamptz,
    ADD COLUMN trial_ends_at          timestamptz,
    ADD COLUMN paid_until             timestamptz,
    ADD COLUMN stripe_customer_id     text UNIQUE,
    ADD COLUMN stripe_subscription_id text;

CREATE TABLE user_tokens (
    id         text PRIMARY KEY,
    user_id    bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind       text NOT NULL CHECK (kind IN ('verify', 'reset')),
    expires_at timestamptz NOT NULL
);
CREATE INDEX user_tokens_user_idx ON user_tokens (user_id);

CREATE TABLE profiles (
    user_id        bigint PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    firstname      text NOT NULL DEFAULT '',
    lastname       text NOT NULL DEFAULT '',
    birthday       date,
    license_number text NOT NULL DEFAULT '',
    created_at     timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE disciplines (
    id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id          bigint REFERENCES users(id) ON DELETE CASCADE,
    key              text UNIQUE,
    name             text NOT NULL,
    series_count     integer NOT NULL CHECK (series_count BETWEEN 1 AND 20),
    shots_per_series integer NOT NULL CHECK (shots_per_series BETWEEN 1 AND 100),
    max_shot         integer NOT NULL DEFAULT 10 CHECK (max_shot BETWEEN 1 AND 100),
    decimal          boolean NOT NULL DEFAULT false,
    target           text NOT NULL DEFAULT '',
    position         integer NOT NULL DEFAULT 100,
    CHECK ((user_id IS NULL) = (key IS NOT NULL)),
    CHECK (NOT decimal OR max_shot = 10)
);
CREATE INDEX disciplines_user_idx ON disciplines (user_id);

INSERT INTO disciplines (key, name, series_count, shots_per_series, decimal, target, position) VALUES
    ('air_pistol', 'Air pistol 10 m', 5, 10, false, 'air_pistol', 1),
    ('air_rifle', 'Air rifle 10 m', 5, 10, true, 'air_rifle', 2),
    ('pistol_25', 'Pistol 25 m', 5, 5, false, 'pistol_precision', 3),
    ('pistol_50', 'Pistol 50 m', 5, 10, false, 'pistol_precision', 4),
    ('rifle_50', 'Rifle 50 m', 5, 10, true, 'rifle_50', 5),
    ('rifle_300', 'Rifle 300 m', 5, 10, false, 'rifle_300', 6);

CREATE TYPE session_kind AS ENUM ('training', 'competition');

CREATE TABLE training_sessions (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    discipline_id bigint NOT NULL REFERENCES disciplines(id),
    kind          session_kind NOT NULL DEFAULT 'training',
    started_at    timestamptz NOT NULL,
    title         text NOT NULL DEFAULT '',
    place         text NOT NULL DEFAULT '',
    notes         text NOT NULL DEFAULT '',
    tags          text[] NOT NULL DEFAULT '{}',
    rating        integer CHECK (rating BETWEEN 1 AND 5),
    client_id     uuid,
    entry_id      bigint REFERENCES entries(id) ON DELETE SET NULL,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, client_id),
    UNIQUE (user_id, entry_id)
);
CREATE INDEX training_sessions_user_idx ON training_sessions (user_id, started_at DESC);
CREATE INDEX training_sessions_discipline_idx ON training_sessions (discipline_id);

CREATE TABLE session_series (
    session_id  bigint NOT NULL REFERENCES training_sessions(id) ON DELETE CASCADE,
    number      integer NOT NULL CHECK (number > 0),
    shots       integer NOT NULL CHECK (shots BETWEEN 1 AND 100),
    total       numeric(7,1) NOT NULL CHECK (total >= 0),
    shot_values numeric(4,1)[],
    xs          real[],
    ys          real[],
    PRIMARY KEY (session_id, number)
);

CREATE TABLE goals (
    user_id       bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    discipline_id bigint NOT NULL REFERENCES disciplines(id) ON DELETE CASCADE,
    target_total  numeric(7,1) NOT NULL CHECK (target_total > 0),
    due_on        date,
    created_at    timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, discipline_id)
);

CREATE TABLE club_plans (
    tenant_id  bigint PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    seats      integer NOT NULL CHECK (seats >= 0),
    paid_until date NOT NULL
);

CREATE TABLE club_seats (
    tenant_id  bigint NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email      text NOT NULL CHECK (email = lower(email)),
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, email)
);
CREATE INDEX club_seats_email_idx ON club_seats (email);
