CREATE TABLE clubs (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id   bigint NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name        text NOT NULL,
    ssv_number  text NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, name),
    UNIQUE (id, tenant_id)
);

CREATE TABLE shooters (
    id             bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id      bigint NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    license_number text,
    firstname      text NOT NULL,
    lastname       text NOT NULL,
    birthday       date,
    postalcode     text NOT NULL DEFAULT '',
    city           text NOT NULL DEFAULT '',
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, license_number),
    UNIQUE (id, tenant_id)
);
CREATE INDEX shooters_name_idx ON shooters (tenant_id, lower(lastname), lower(firstname));

CREATE TABLE shooter_clubs (
    tenant_id  bigint NOT NULL,
    shooter_id bigint NOT NULL,
    club_id    bigint NOT NULL,
    PRIMARY KEY (shooter_id, club_id),
    FOREIGN KEY (shooter_id, tenant_id) REFERENCES shooters(id, tenant_id) ON DELETE CASCADE,
    FOREIGN KEY (club_id, tenant_id) REFERENCES clubs(id, tenant_id) ON DELETE CASCADE
);
CREATE INDEX shooter_clubs_club_idx ON shooter_clubs (club_id);
