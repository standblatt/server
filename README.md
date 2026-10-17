# Standblatt

Competition management for shooting ranges (Swiss SSV rules): registration, lane assignment, result entry, awards, rankings, club competition and accounting. A single Go server with an embedded SvelteKit app; PostgreSQL is the only dependency. Website: [www.standblatt.ch](https://www.standblatt.ch).

> Standblatt is under active development; the first clubs use it from 2027. Expect changes until 1.0.

**Documentation:** the help is part of the app under `/help` (sources in `web/src/help/`); operations and development docs are in `docs/`.

## Development

Requirements: Go, Node, Docker (e2e needs Docker only).

```sh
make dev      # PostgreSQL, Mailpit, backend and Vite in the foreground, app on http://localhost:5173; Ctrl-C stops everything
make prepare  # before a commit: OpenAPI spec and TS types, PO catalogs, gofmt, go vet, doc comments, svelte-check
make test     # Go unit and integration tests against the compose PostgreSQL, then the Playwright e2e specs
make docs     # help: update the PO catalogs, write the translated pages, regenerate the screenshots
```

`make dev` points git at the hooks in `.githooks`, whose pre-commit hook runs `make prepare` and stops the commit when it changed files (`git commit --no-verify` skips it).

`make dev` seeds a local login `dev` / `devpassword` (club "Dev Club" with a range and an event). The admin panel is at http://localhost:5173/admin with the token from `.env`.

## Configuration

Environment variables only. For development copy `.env.example` to `.env` (`make dev` does this on first start).

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | required | PostgreSQL URL |
| `STANDBLATT_LISTEN` | `:8080` | Listen address |
| `STANDBLATT_METRICS_LISTEN` | `:9090` | Prometheus `/metrics` on a separate listener on purpose, so it is never reachable through the ingress. Empty: no metrics |
| `STANDBLATT_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error`; JSON lines on stderr |
| `STANDBLATT_TRUST_PROXY` | `false` | Take the client IP from `X-Forwarded-For` (only behind your own ingress) |
| `STANDBLATT_ADMIN_TOKEN` | empty | Token for the admin panel at `/admin`, at least 32 characters. Empty: no admin panel |
| `STANDBLATT_PUBLIC_URL` | empty | Base of links in mails; empty: scheme and host of the request |
| `STANDBLATT_SMTP_HOST` | empty | Mail server; empty: mails are only logged |
| `STANDBLATT_SMTP_PORT` | `587` | `465` uses implicit TLS, other ports STARTTLS when the server offers it |
| `STANDBLATT_SMTP_USER`, `STANDBLATT_SMTP_PASSWORD` | empty | Login at the mail server, if it needs one |
| `STANDBLATT_SMTP_FROM` | `Standblatt <noreply@standblatt.ch>` | Sender of mails |
| `STANDBLATT_SIGNUP` | `false` | Open sign-up for the personal training log |
| `STANDBLATT_ENVIRONMENT` | | Deployment name such as `staging`; anything but empty or `production` shows a red banner and a title suffix in the app |
| `STANDBLATT_TOTP` | `true` | Users may add an authenticator app as second factor |
| `STANDBLATT_PRICE_LABEL` | empty | Price shown to users, display only |
| `STANDBLATT_MATOMO_URL`, `STANDBLATT_MATOMO_SITE_ID` | empty | Matomo web analytics without cookies; both or neither |

Migrations are embedded and applied at startup.

## Accounts and clubs

Club accounts are created in the admin panel; shooters can also sign up themselves for the personal training log at `/my` when `STANDBLATT_SIGNUP=true` (see `docs/content/my-training/`). The admin panel (`/admin`, needs `STANDBLATT_ADMIN_TOKEN`) creates clubs and users and assigns users to clubs as `admin` or `staff`. Club admins then manage their ranges, events and members in the app. Staff can use events but not change club settings.

## Shooter master data

Every club keeps its own shooters and clubs; nothing is shared between tenants. An import of official SSV data is not available (missing). Club admins import CSV files on the Shooters page:

```
license;firstname;lastname;birthday;postalcode;city;club
100100;Sandra;Example;1969-06-02;9050;Appenzell;PS Example
```

`firstname` and `lastname` are required, the other columns optional; `;` or `,` as separator, birthdays as `YYYY-MM-DD` or `DD.MM.YYYY`. One line per shooter and club. Shooters are matched by license number (or by name and birthday when there is none), purely numeric license numbers are padded to six digits, missing clubs are created and nothing is deleted. A file with invalid lines is rejected as a whole with a list of the lines to fix.

## Programmes

Each event has its own programmes: series (shots, maximum points, mouches), prices per age category (only categories with a price may shoot the programme) and award tiers with a minimum total per age category. Programmes with the same name are variants; a shooter registers for one of them. A programme can be marked as combinable with another programme of the event.

Programmes come from the editor, from an earlier event of the club (copy), from a programme file (`*.programme.json`, made by Export) or from the global library, which the platform admin maintains in the admin panel by uploading programme files.

## Registration and live updates

The event dashboard lists waiting and shooting shooters and the lanes. `F4` focuses the search; a scanned registration code (`REGCODEAAAA<registration>AAAAL<license or 0>`) followed by Enter opens that registration. A lane can hold one active registration at a time and a shooter can register for each programme of an event once, both enforced by the database.

Results are entered per series on the result page of a registration (the scanned code opens it). The cursor moves on as soon as another digit would exceed the series maximum, Enter saves. Combined programmes copy the chosen series of their source programme; once every programme of a registration has results, the registration completes and frees its lane. Only awards the total reaches in the shooter's age category can be chosen.

Rankings are computed in one place (`internal/domain/ranking.go`): higher total; then the series compared from the last one backwards; then more mouches; then age (among U10–U21 the younger, among the others the older shooter first, non-juniors before juniors). Ranks are consecutive.

For public screens a club admin creates a kiosk link in the event setup. `/kiosk/<token>` shows the live rankings (and the lane status if the event enables it) without login; the token only grants read access to that event and can be revoked at any time.

Browsers keep a WebSocket to `/api/tenants/{tenant}/events/{event}/live`. Changes are published with PostgreSQL `NOTIFY`; every replica listens and forwards them to its own clients, so no sticky sessions or extra broker are needed. The ingress must allow WebSocket upgrades (Traefik does by default).

## Lane reservations

Staff reserve lanes of a range for a time (Reservations tab of an event), e.g. lanes `1, 5, 7-10` from 18:00 to 20:00, shown as reserved 30 minutes before. Reservations may not overlap on a lane; back-to-back is fine. Lanes show as out of service, occupied, reserved or free, in that priority; a reserved lane can still be assigned after a confirmation. A ticker in the server (one replica at a time, via a PostgreSQL advisory lock) announces when reservations start showing or end, so dashboards and kiosks update without polling.

## Statistics, club competition and exports

The Statistics tab of an event shows registrations, completion, revenue, participation, age categories and awards, and the club competition (only programmes marked "counts for the club competition"; the first two series of each completed result count; a club needs five participants; of n participants the best max(5, ⌊n/2⌋) are mandatory, the rest count 2 %; the result is truncated to three decimals). Every report is available as PDF and as CSV for spreadsheets (UTF-8 with BOM, `;` separated).

## Documents and printing

Lane sheet (A5), control strip (label 90 × 29 mm), result sheet (A5), wreath card (A6, printed onto pre-printed cards) and rankings (A4) are rendered as PDF by the server (`internal/docs`, go-pdf/fpdf). Every document can be opened as PDF.

Printing at the range goes through the print agent (`cmd/print-agent`), a small static binary for a computer next to the printers (e.g. a Raspberry Pi with CUPS). It only makes outgoing HTTPS requests: it reports the CUPS printers with their options, waits for jobs (long poll) and prints them with `lp`.

Installation and configuration: [docs/admin/print-agent.md](docs/admin/print-agent.md). In the app, pick a printer and its options (tray, paper size, media type, further `lp` options) for each document.

When no printer is set up for a document or no agent of the range is online, the Print buttons open the PDF instead. Jobs whose agent does not confirm within two minutes are retried, after three attempts they fail; finished jobs are deleted after a week.

## Deployment

Helm chart in `helm/`, see [docs/admin/kubernetes.md](docs/admin/kubernetes.md); image `ghcr.io/standblatt/server`, chart `oci://ghcr.io/standblatt/charts/standblatt`, both built by GitHub Actions. The hosted installation is deployed by Flux and Keel. `DATABASE_URL` and `STANDBLATT_ADMIN_TOKEN` come from `envFromSecret`.
