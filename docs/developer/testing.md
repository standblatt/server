# Testing

## Go

```sh
make test   # Go tests, then e2e; for a quick loop: go test ./internal/domain/
```

- `internal/domain`: pure unit tests (ranking, club competition, special prizes, age categories, amounts in words …).
- `internal/store` and `internal/web`: integration tests. `internal/testdb` creates a throwaway, migrated database per test from `TEST_DATABASE_URL`; `make test` starts the compose PostgreSQL on port 5433 and sets the variable. Without it these tests are skipped.
- `internal/web` runs twice (`TestMain` in `main_test.go`): once with sessions and once with API tokens, where every logged-in test client sends a token with every scope of its user instead of the session cookie on each path a token reaches. After a full run, `TestMain` fails if an operation was never called successfully, or a token-reachable one never with a token: a new endpoint needs a test that calls it. `TestPermissionTable` and `TestPermissionLeaks` check who may call each club operation (see [modules](modules.md#permissions)).
- `internal/docs`: checks the text of rendered PDFs with `pdftotext` when installed. `DOCS_OUT=dir go test ./internal/docs/` writes the PDFs for a visual check.
- `internal/i18n`: fails on untranslated messages or mismatched format verbs.

## End to end

`make test` runs them after the Go tests. It builds the real image, starts it with an empty database on `:3002` (compose profile `e2e`) and runs the Playwright specs in `e2e/tests` in the `mcr.microsoft.com/playwright` image with Chromium at 1280×800, several spec files in parallel. `@playwright/test` in `e2e/package.json` is the only place of the Playwright version: the Makefile takes the image tag from it and CI installs the browsers of that version, so a Dependabot update of the package moves both.

On a failure, `e2e/test-results` holds a screenshot and a trace of the failed test; open the trace with `npx playwright show-trace <file>` in `e2e/`. The CI job keeps the folder and the server log as an artifact when it fails.

Specs create their own club through the admin panel (`createClub` in `e2e/helpers.js`) and their data through the API, so they can run in parallel against the same instance. Each file runs serially: tests in a file may build on each other, and a retry repeats the whole file with a new club. `c(scope, selector, text)` finds the innermost element containing a text, like `cy.contains`. Toasts are `[data-sonner-toast]`, confirmations `[role=alertdialog]`.

## Documentation screenshots

`make docs` runs `e2e/docs/screenshots.spec.js` as the Playwright project `docs` (1440×900) and writes the PNGs to `web/static/help/screenshots`. It is not part of `make test`. Use invented names only: the screenshots are public.
