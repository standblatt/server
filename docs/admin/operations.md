# Operations

## Logs

Standblatt writes JSON lines to stderr: startup, background work, warnings and errors (requests are not logged one by one). `STANDBLATT_LOG_LEVEL=debug` adds more detail.

## Metrics

Prometheus metrics are served on `STANDBLATT_METRICS_LISTEN` (default `:9090`) at `/metrics`, on a separate port so they never leak through the ingress. The Helm chart can create a `PodMonitor`.

## Backups

All state is in PostgreSQL, including print jobs, sessions and realtime events. Back up the database (for example with the CloudNativePG backup of the cluster); the containers hold nothing worth saving.

## Updates

Deploy a newer image. Migrations run at startup, so a new version migrates the database before it serves requests. Read the [changelog](../../CHANGELOG.md) before upgrading.

Open browser tabs check for a new version every minute. They then show a *Reload* notice and load the new version on the next page change at the latest; kiosk displays reload by themselves.

## Background work

Every replica runs small janitors (expired sessions, old print jobs, requeueing unconfirmed jobs). The reservation ticker announces reservation starts and ends every 30 seconds; a PostgreSQL advisory lock makes sure only one replica does it at a time.
