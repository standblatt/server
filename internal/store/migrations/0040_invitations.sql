-- Invitations by email: the account and the membership are made when the invitation is accepted through its link.
CREATE TABLE invitations (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id   bigint NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email       text NOT NULL,
    role_id     bigint NOT NULL,
    token_hash  text NOT NULL UNIQUE,
    invited_by  bigint REFERENCES users(id) ON DELETE SET NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    expires_at  timestamptz NOT NULL,
    accepted_at timestamptz,
    FOREIGN KEY (role_id, tenant_id) REFERENCES roles(id, tenant_id) ON DELETE CASCADE
);
CREATE UNIQUE INDEX invitations_open ON invitations (tenant_id, lower(email)) WHERE accepted_at IS NULL;
