# Development

Requirements: Go, Node and Docker.

```sh
make dev      # PostgreSQL, Mailpit, backend and Vite in the foreground, app on http://localhost:5173; Ctrl-C stops everything
make prepare  # before a commit: OpenAPI spec and TS types, PO catalogs, gofmt, go vet, staticcheck, doc comments, svelte-check
make test     # Go unit and integration tests against the compose PostgreSQL, then the Playwright e2e specs
make docs     # help: update the PO catalogs, write the translated pages, regenerate the screenshots
make clean    # stop the containers, remove node_modules and build output (keeps .env and the database volume)
```

`make dev` also sets `git config core.hooksPath .githooks`: the pre-commit hook runs `make prepare` and stops the commit when it changed files (new strings in the catalogs, a regenerated `api.gen.ts`); review and stage them, then commit again. `git commit --no-verify` skips it. The image is built by `docker build -f docker/Dockerfile .`; the print agent is built in [standblatt/print-agent](https://github.com/standblatt/print-agent).

`make dev` seeds the login `dev` / `devpassword` with a club, a range and an event. The admin panel is at `/admin` with the token from `.env` (created from `.env.example` on the first start).

## Conventions

- Everything in the repository is English; translations live in the PO catalogs only.
- No comments inside function bodies unless really necessary, but a doc comment on every function: godoc (`// name …`) on every Go function and method, a one-line JSDoc (`/** … */`) on every TypeScript and Svelte function. `make prepare` (via `scripts/funcdoc.py`, also a CI job) fails on missing ones; the shadcn copies in `components/ui`, `api.gen.ts` and the locales are exempt.
- Rules belong in `internal/domain` with table tests, not in handlers or the frontend.
- huma path parameters must be at most one embedding level deep in input structs (use the flat `*Path` structs); arrays are non-nullable, so store list methods never return nil slices.
- UI components come from shadcn-svelte (`web/src/lib/components/ui`), no other UI kits. `ask()` instead of `confirm()`, toasts for feedback, dialogs and sheets for create and edit.
- Every change gets an entry under `[Unreleased]` in `CHANGELOG.md`.
