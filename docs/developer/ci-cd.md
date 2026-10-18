# CI/CD

GitHub Actions (`.github/workflows/ci.yml`) builds and tests on every push to master, every tag `vX.Y.Z` and every pull request; it does not deploy.

| Jobs | What |
|---|---|
| `go`, `govulncheck`, `funcdoc`, `helm`, `web-check`, `web` | gofmt, vet, staticcheck and tests with PostgreSQL; known vulnerabilities in called code (does not block images); every function has a doc comment; chart lint; svelte-check; build of the SPA (`binary` waits only for this one) |
| `binary`, `print-agent` | the server binary (built once, with `git describe` as version) and the print agent for linux amd64/arm64 |
| `e2e` | Playwright against the binary, with PostgreSQL and Mailpit as services |
| `docker-image` | pushes and tags the image on ghcr.io (master and tags only) |
| `chart`, `release` | tags only: the chart as OCI artifact to `oci://ghcr.io/standblatt/charts`, a GitHub release with the changelog section and the print agent binaries |

## Deployments

| Trigger | Image tags | App |
|---|---|---|
| push to master | `<sha>`, `bleeding_edge` | app.staging.standblatt.ch |
| tag `vX.Y.Z` | `<sha>`, `latest`, `vX.Y.Z` | app.standblatt.ch |

`docker-image` pushes the image once `e2e` passed. Flux in the cluster's GitOps repo installs the chart that `chart` pushes on every tag (an `OCIRepository` following new chart versions, a HelmRelease per environment); Keel polls the moving image tag (staging every minute, production every five minutes) and rolls out a new digest. A chart change goes live with the next tag. Image and chart are public, so neither the workflow nor the cluster needs credentials beyond `GITHUB_TOKEN`.

The image version is `git describe --tags --always` (e.g. `v0.3.0` or `v0.3.0-2-g1a2b3c4`); the sidebar of the app shows it.

## Notes

- `docker-image` does not compile anything: it packs the binary from the job `binary`, the same one `e2e` tested, into a distroless image (`docker/ci.Dockerfile`). Master and tag runs of one commit reuse the image of the SHA tag, so both serve the same SPA build. `docker/Dockerfile` builds everything from source (`docker build -f docker/Dockerfile .`).
- Go and npm downloads are cached by `actions/setup-go` and `actions/setup-node`.
- A newer push to the same branch cancels the running workflow; tag runs always finish.
- CodeQL runs as GitHub's *default setup* (no workflow file) for Go, JavaScript/TypeScript and the workflows; findings are under *Security → Code scanning*. Dependabot security alerts are on.
- Dependabot (`.github/dependabot.yml`) proposes weekly updates for Go modules, both npm projects, the workflow actions and the Docker images; the pinned major versions (Go 1.27, Node 26, PostgreSQL 18, SvelteKit 2, TypeScript 6) are ignore rules there.
