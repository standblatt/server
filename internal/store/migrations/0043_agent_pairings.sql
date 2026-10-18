-- Pairing by code: the desktop print agent asks for a code, a member confirms it on /connect and the agent fetches
-- its token once with its poll secret. No agent token is kept here: confirming creates the agent with a placeholder
-- hash and fetching sets the real one.
CREATE TABLE agent_pairings (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code        text NOT NULL UNIQUE,
    secret_hash text NOT NULL UNIQUE,
    name        text NOT NULL,
    os          text NOT NULL,
    version     text NOT NULL,
    expires_at  timestamptz NOT NULL,
    tenant_id   bigint REFERENCES tenants ON DELETE CASCADE,
    agent_id    bigint REFERENCES print_agents ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now()
);
