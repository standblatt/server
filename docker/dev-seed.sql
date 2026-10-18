-- Local development data, applied by `make dev`. Login: dev / devpassword
INSERT INTO users (name, password_hash)
VALUES ('dev', '$2a$12$GgZO./3wHssJNElV.flSRuwmqxTOyrbGamaBCgSwDlxsTUyVfhG5y')
ON CONFLICT (name) DO NOTHING;

INSERT INTO tenants (name) VALUES ('Dev Club') ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (tenant_id, key, permissions, position)
SELECT t.id, d.key, d.permissions, d.position FROM tenants t, (VALUES
    ('admin', '{}'::text[], 1),
    ('board', '{events.view,events.operate,events.configure,reports.view,shooters.view,shooters.edit,shooters.manage,forms.manage}', 2),
    ('staff', '{events.view,events.operate,reports.view,shooters.view,shooters.edit}', 3),
    ('member', '{events.view}', 4)
) d(key, permissions, position) WHERE t.name = 'Dev Club'
ON CONFLICT DO NOTHING;

INSERT INTO memberships (tenant_id, user_id, role_id)
SELECT t.id, u.id, r.id FROM tenants t JOIN roles r ON r.tenant_id = t.id AND r.key = 'admin', users u
WHERE t.name = 'Dev Club' AND u.name = 'dev'
ON CONFLICT DO NOTHING;

INSERT INTO ranges (tenant_id, name, lane_count)
SELECT id, 'Main range', 10 FROM tenants WHERE name = 'Dev Club'
ON CONFLICT DO NOTHING;

INSERT INTO events (tenant_id, range_id, name, starts_on)
SELECT t.id, s.id, 'Test event', current_date FROM tenants t JOIN ranges s ON s.tenant_id = t.id
WHERE t.name = 'Dev Club' AND s.name = 'Main range'
  AND NOT EXISTS (SELECT 1 FROM events e WHERE e.tenant_id = t.id);
