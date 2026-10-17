CREATE TABLE programmes (
    id                          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id                   bigint NOT NULL,
    event_id                    bigint NOT NULL,
    name                        text NOT NULL,
    variant                     text NOT NULL DEFAULT '',
    description                 text NOT NULL DEFAULT '',
    barcode                     text NOT NULL DEFAULT '',
    counts_for_club_competition boolean NOT NULL DEFAULT false,
    show_on_rankings            boolean NOT NULL DEFAULT true,
    active                      boolean NOT NULL DEFAULT true,
    combinable_with_id          bigint REFERENCES programmes(id) ON DELETE SET NULL,
    created_at                  timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY (event_id, tenant_id) REFERENCES events(id, tenant_id) ON DELETE CASCADE,
    UNIQUE (event_id, name, variant),
    UNIQUE (id, tenant_id),
    CHECK (combinable_with_id <> id)
);

CREATE TABLE programme_series (
    programme_id  bigint NOT NULL REFERENCES programmes(id) ON DELETE CASCADE,
    number        integer NOT NULL CHECK (number > 0),
    shots         integer NOT NULL CHECK (shots > 0),
    max_points    integer NOT NULL CHECK (max_points > 0),
    name          text NOT NULL DEFAULT '',
    track_mouches boolean NOT NULL DEFAULT false,
    PRIMARY KEY (programme_id, number)
);

CREATE TABLE programme_prices (
    programme_id bigint NOT NULL REFERENCES programmes(id) ON DELETE CASCADE,
    age_category age_category NOT NULL,
    price        numeric(8,2) NOT NULL CHECK (price >= 0),
    PRIMARY KEY (programme_id, age_category)
);

CREATE TYPE award_kind AS ENUM ('natural', 'card');

CREATE TABLE award_tiers (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    programme_id bigint NOT NULL REFERENCES programmes(id) ON DELETE CASCADE,
    position     integer NOT NULL,
    name         text NOT NULL,
    kind         award_kind NOT NULL,
    value        numeric(8,2) NOT NULL CHECK (value >= 0),
    UNIQUE (programme_id, position)
);

CREATE TABLE award_limits (
    tier_id      bigint NOT NULL REFERENCES award_tiers(id) ON DELETE CASCADE,
    age_category age_category NOT NULL,
    points       numeric(6,2) NOT NULL CHECK (points >= 0),
    PRIMARY KEY (tier_id, age_category)
);

CREATE TABLE programme_templates (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       text NOT NULL,
    variant    text NOT NULL DEFAULT '',
    definition jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (name, variant)
);
