# API

The JSON API lives under `/api` and is written with [huma](https://huma.rocks/), code first. The server publishes:

- `/api/openapi` (OpenAPI 3.1 spec, also `standblatt openapi` on the command line),
- `/api/docs` (interactive documentation).

The frontend uses `openapi-fetch` with types generated from the spec: after changing an operation, run `make prepare` (the pre-commit hook does it too) and commit `web/openapi.json` and `web/src/lib/api.gen.ts`.

## Authentication

| Client | How |
|---|---|
| App | session cookie `standblatt_session` from `POST /api/auth/login` |
| Kiosk screen | `Authorization: Bearer <kiosk token>` on `/api/kiosk/*`, `?token=` on the WebSocket |
| Print agent | `Authorization: Bearer <agent token>` on `/api/agent/*` |
| Scripts and tools | personal API token `sb_…` as `Authorization: Bearer` or `X-API-Key` on `/api/tenants/*`, the training paths of `/api/my` and `GET /api/token` |

### Personal API tokens

Created in the account (`POST /api/me/tokens`, session only) with a name, scopes, the clubs (`tenants`) and an optional expiry of 30, 90 or 365 days. Only the sha256 hash is stored; the secret (`sb_` + 40 hex characters) is returned once.

| Scope | Allows |
|---|---|
| `events.view` … `club.manage` | the operations of `/api/tenants/{tenantId}/…` that need this permission (the table is `internal/web/testdata/permissions.txt`) |
| `training:read` | `GET` on `/api/my` and `/api/my/{trainings,stats,disciplines,weapons,goals,competitions,export}` |
| `training:write` | `POST`, `PUT` and `DELETE` on those paths |

A token may use a club's operations only for the clubs it was created for, and only with the permissions that are both its scopes and in the owner's current role there: a club it is not for answers 404, a missing permission 403. Creating a token refuses clubs the owner is not in and club scopes the owner's role has in none of them (422). Account endpoints (`/api/me*`, `/api/auth/*`), the token management and everything else answer 403. Unknown, revoked and expired tokens get 401 with a distinct `detail`. `GET /api/token` returns the token's name, owner, scopes, clubs and expiry.

Limits: a token makes at most `STANDBLATT_API_RATE_LIMIT` requests per five minutes (default 600), more get 429 with `Retry-After`; an address that sends 20 unknown tokens within 15 minutes gets 429 too. Each use updates *last used*.

```sh
curl -H "Authorization: Bearer $TOKEN" https://app.standblatt.ch/api/token
curl -H "Authorization: Bearer $TOKEN" https://app.standblatt.ch/api/tenants/12/events/34/rankings
```

## Conventions

- Tenant-scoped paths: `/api/tenants/{tenantId}/…`; non-members get 404, members whose role lacks the permission of the operation get 403 (see `internal/permissions`).
- Errors are RFC 9457 problem details (`application/problem+json`) with an English `detail`.
- A change to an archived event answers 409.
- WebSockets: `/api/tenants/{tenantId}/events/{eventId}/live` and `/api/kiosk/live` send `{topic}` messages; clients refetch what the topic covers.
