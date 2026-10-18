-- API tokens reach club data: scopes are permission keys (internal/permissions) plus training:read/write, tenants are
-- the clubs a token may use. Tokens of the training-only kind are not carried over; owners create new ones.
DELETE FROM api_tokens;
ALTER TABLE api_tokens ADD COLUMN tenants bigint[] NOT NULL DEFAULT '{}';
