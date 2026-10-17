CREATE TABLE print_agents (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id    bigint NOT NULL,
    range_id     bigint NOT NULL,
    name         text NOT NULL,
    token_hash   text NOT NULL UNIQUE,
    printers     jsonb NOT NULL DEFAULT '[]',
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz,
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE
);

CREATE TABLE print_settings (
    tenant_id bigint NOT NULL,
    range_id  bigint NOT NULL,
    document  text NOT NULL CHECK (document IN ('lane_sheet', 'control_strip', 'result_sheet', 'wreath_card', 'ranking')),
    printer   text NOT NULL,
    options   jsonb NOT NULL DEFAULT '{}',
    PRIMARY KEY (range_id, document),
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE
);

CREATE TYPE print_status AS ENUM ('queued', 'printing', 'done', 'failed');

CREATE TABLE print_jobs (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id   bigint NOT NULL,
    range_id    bigint NOT NULL,
    event_id    bigint,
    document    text NOT NULL,
    title       text NOT NULL,
    pdf         bytea NOT NULL,
    printer     text NOT NULL,
    options     jsonb NOT NULL DEFAULT '{}',
    status      print_status NOT NULL DEFAULT 'queued',
    error       text NOT NULL DEFAULT '',
    attempts    integer NOT NULL DEFAULT 0,
    created_at  timestamptz NOT NULL DEFAULT now(),
    claimed_at  timestamptz,
    finished_at timestamptz,
    FOREIGN KEY (range_id, tenant_id) REFERENCES ranges(id, tenant_id) ON DELETE CASCADE
);
CREATE INDEX print_jobs_queue_idx ON print_jobs (range_id, created_at) WHERE status = 'queued';
CREATE INDEX print_jobs_range_idx ON print_jobs (range_id, created_at DESC);
