# Configuration

Standblatt reads environment variables only.

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | required | PostgreSQL URL, e.g. `postgres://user:pass@host:5432/db?sslmode=require` |
| `STANDBLATT_LISTEN` | `:8080` | listen address of the app |
| `STANDBLATT_METRICS_LISTEN` | `:9090` | Prometheus `/metrics` on a separate listener, so it is never reachable through the ingress. Empty: no metrics |
| `STANDBLATT_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error`; JSON lines on stderr |
| `STANDBLATT_TRUST_PROXY` | `false` | take the client address from `X-Forwarded-For` (only behind your own ingress); used for login rate limits |
| `STANDBLATT_ADMIN_TOKEN` | empty | token for the [admin panel](admin-panel.md), at least 32 characters. Empty: the admin panel is disabled |
| `STANDBLATT_PUBLIC_URL` | empty | base of the links in mails, e.g. `https://app.standblatt.ch`. Empty: scheme and host of the request |
| `STANDBLATT_SMTP_HOST` | empty | mail server. Empty: mails are only logged |
| `STANDBLATT_SMTP_PORT` | `587` | `465` uses implicit TLS, other ports STARTTLS when the server offers it |
| `STANDBLATT_SMTP_USER`, `STANDBLATT_SMTP_PASSWORD` | empty | login at the mail server, if it needs one |
| `STANDBLATT_SMTP_FROM` | `Standblatt <noreply@standblatt.ch>` | sender of mails |
| `STANDBLATT_SIGNUP` | `false` | anyone can create an account for [My training](../../web/src/help/en/my-training/index.md) |
| `STANDBLATT_ENVIRONMENT` | | deployment name such as `staging`; anything but empty or `production` shows a red *STAGING* banner and adds the name to the page title, so nobody confuses test and live |
| `STANDBLATT_TOTP` | `true` | users may add an authenticator app as second factor (see [Login and language](../../web/src/help/en/user-guide/login-language.md)); `false` hides the option and ignores stored secrets |
| `STANDBLATT_MATOMO_URL`, `STANDBLATT_MATOMO_SITE_ID` | empty | web analytics with [Matomo](https://matomo.org): base URL of the instance (e.g. `https://analytics.example.org/`) and the numeric site ID; set both or neither. The app reports page views without cookies, with record IDs replaced by `:id` and without query strings; kiosk displays are not counted |

Generate the admin token with `openssl rand -base64 32` and keep it in a secret store.

> **Stateless:** Sessions, login rate limits, maintenance mode, print jobs and realtime events all live in PostgreSQL. Any number of replicas can run behind a load balancer without sticky sessions.
