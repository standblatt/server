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
| Scripts of a shooter | personal API token `sb_…` as `Authorization: Bearer` or `X-API-Key` on `/api/my/*` |

### Personal API tokens

Created in My training → Settings (`POST /api/my/tokens`) with a name, scopes and an optional expiry of 30, 90 or 365 days. Only the sha256 hash is stored; the secret (`sb_` + 40 hex characters) is returned once. Scopes:

| Scope | Allows |
|---|---|
| `training:read` | `GET` on the paths below |
| `training:write` | `POST`, `PUT` and `DELETE` on the paths below |

A token reaches `/api/my` and `/api/my/{trainings,stats,disciplines,weapons,goals,competitions,export}`. Everything else answers 403: profile, email, account, the token management itself; a token on any other path is ignored (401). Unknown, revoked and expired tokens get 401 with a distinct `detail`; a missing scope gets 403. Each use updates *last used*.

```sh
curl -H "Authorization: Bearer $TOKEN" https://app.standblatt.ch/api/my/trainings?limit=10
```

## Conventions

- Tenant-scoped paths: `/api/tenants/{tenantId}/…`; non-members get 404, members whose role lacks the permission of the operation get 403 (see `internal/permissions`).
- Errors are RFC 9457 problem details (`application/problem+json`) with an English `detail`.
- A change to an archived event answers 409.
- WebSockets: `/api/tenants/{tenantId}/events/{eventId}/live` and `/api/kiosk/live` send `{topic}` messages; clients refetch what the topic covers.
