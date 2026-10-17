# Installation

Standblatt is a single Go binary with the web app embedded, shipped as a container image. PostgreSQL is the only dependency.

## Requirements

- PostgreSQL 18 (any 15+ should work), one database per installation.
- The container image `ghcr.io/standblatt/server` (distroless, runs as a non-root user, read-only root filesystem works).
- A reverse proxy or ingress for HTTPS that passes WebSocket upgrades (Traefik does by default).

## Running the container

```sh
docker run -d --name standblatt -p 8080:8080 \
  -e DATABASE_URL='postgres://standblatt:secret@db:5432/standblatt?sslmode=require' \
  -e STANDBLATT_ADMIN_TOKEN="$(openssl rand -base64 32)" \
  -e STANDBLATT_TRUST_PROXY=true \
  ghcr.io/standblatt/server:latest
```

- Port **8080** serves the app, the API and the admin panel; port **9090** serves Prometheus metrics (see [Operations](operations.md)).
- The database schema is migrated at startup. Several replicas may start at the same time; an advisory lock makes them wait for each other.
- All configuration comes from environment variables, see [Configuration](configuration.md).

Image tags: the short commit SHA of every build, `bleeding_edge` for the latest master build, `latest` for the latest release and `vX.Y.Z` for every release. Pin a release tag in production.

## Health checks

| Path | Meaning |
|---|---|
| `/healthz` | the process is up (liveness) |
| `/readyz` | the database answers within two seconds (readiness) |

## First steps

1. Open `/admin` and log in with `STANDBLATT_ADMIN_TOKEN`.
2. Create a club and a user, and make the user admin of the club (see [Admin panel](admin-panel.md)).
3. The user logs in to the app and continues with [Your first event](../../web/src/help/en/getting-started/first-event.md).

For Kubernetes, use the Helm chart, see [Kubernetes](kubernetes.md).
