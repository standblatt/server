-- Roles per club made of the permissions in internal/permissions; memberships point to a role instead of admin/staff.
-- key marks the roles every club starts with (name NULL = the translated default name); the admin role has every permission.
CREATE TABLE roles (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id   bigint NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    key         text,
    name        text,
    permissions text[] NOT NULL DEFAULT '{}',
    position    integer NOT NULL DEFAULT 0,
    UNIQUE (id, tenant_id),
    UNIQUE (tenant_id, key)
);

INSERT INTO roles (tenant_id, key, permissions, position)
SELECT t.id, d.key, d.permissions, d.position FROM tenants t, (VALUES
    ('admin', '{}'::text[], 1),
    ('board', '{events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,forms.manage}', 2),
    ('staff', '{events.view,events.operate,reports.view,shooters.view,shooters.edit}', 3),
    ('member', '{events.view}', 4)
) d(key, permissions, position);

ALTER TABLE memberships ADD COLUMN role_id bigint;
UPDATE memberships m SET role_id = r.id FROM roles r WHERE r.tenant_id = m.tenant_id AND r.key = m.role::text;
ALTER TABLE memberships ALTER COLUMN role_id SET NOT NULL;
ALTER TABLE memberships ADD FOREIGN KEY (role_id, tenant_id) REFERENCES roles(id, tenant_id);
ALTER TABLE memberships DROP COLUMN role;
DROP TYPE tenant_role;
CREATE INDEX memberships_role_idx ON memberships (role_id);
