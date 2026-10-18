# Standblatt

Competition management for shooting clubs (Swiss SSV rules): registration at the desk, lanes and reservations, results, awards, rankings, kiosk displays and reports. Every event switches on the modules it needs, from a shooting festival with online booking to an annual programme without a desk. Shooters keep a personal training log in *My training*.

A single Go server with an embedded SvelteKit app, in English, German, French and Italian; PostgreSQL is the only dependency. Website: [www.standblatt.ch](https://www.standblatt.ch).

> Standblatt is under active development; the first clubs use it from 2027. Expect changes until 1.0.

## Documentation

- **Help** for club admins, staff and shooters: part of the app under `/help` (sources in `web/src/help/`).
- **Operations and development**: [docs/](docs/README.md), e.g. [installation](docs/admin/installation.md), [configuration](docs/admin/configuration.md), [Kubernetes](docs/admin/kubernetes.md), [print agent](docs/admin/print-agent.md), [architecture](docs/developer/architecture.md) and [modules](docs/developer/modules.md).
- **Changes**: [CHANGELOG.md](CHANGELOG.md).

## Development

Requirements: Go, Node and Docker.

```sh
make dev      # PostgreSQL, Mailpit, backend and Vite; app on http://localhost:5173, login dev / devpassword
make prepare  # before a commit (also the pre-commit hook): API types, catalogs, gofmt, vet, staticcheck, svelte-check
make test     # Go tests, then the Playwright e2e specs
make docs     # help translations and screenshots
make clean    # containers down, node_modules and build output removed
```

More in [docs/developer/development.md](docs/developer/development.md).

## Deployment

The image `ghcr.io/standblatt/server` and the Helm chart `oci://ghcr.io/standblatt/charts/standblatt` are built by GitHub Actions; configuration is environment variables only and migrations run at startup. Printing at the range goes through the small print agent (`cmd/print-agent`), released with every version.

## License

[AGPL-3.0](LICENSE)
