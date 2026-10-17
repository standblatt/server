-- Local development data, applied by `make dev`. Login: dev / devpassword
INSERT INTO users (name, password_hash)
VALUES ('dev', '$2a$12$GgZO./3wHssJNElV.flSRuwmqxTOyrbGamaBCgSwDlxsTUyVfhG5y')
ON CONFLICT (name) DO NOTHING;

INSERT INTO tenants (name) VALUES ('Dev Club') ON CONFLICT (name) DO NOTHING;

INSERT INTO memberships (tenant_id, user_id, role)
SELECT t.id, u.id, 'admin' FROM tenants t, users u WHERE t.name = 'Dev Club' AND u.name = 'dev'
ON CONFLICT DO NOTHING;

INSERT INTO ranges (tenant_id, name, lane_count)
SELECT id, 'Main range', 10 FROM tenants WHERE name = 'Dev Club'
ON CONFLICT DO NOTHING;

INSERT INTO events (tenant_id, range_id, name, starts_on)
SELECT t.id, s.id, 'Test event', current_date FROM tenants t JOIN ranges s ON s.tenant_id = t.id
WHERE t.name = 'Dev Club' AND s.name = 'Main range'
  AND NOT EXISTS (SELECT 1 FROM events e WHERE e.tenant_id = t.id);
