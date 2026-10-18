ifneq (,$(wildcard .env))
include .env
export $(shell sed -E 's/^[[:space:]]*\#.*//; s/=.*//' .env)
endif

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null)
COMPOSE := docker compose -f $(CURDIR)/docker/compose.yaml
# The image follows @playwright/test in e2e/package.json, the only place Dependabot updates.
PLAYWRIGHT_IMAGE := mcr.microsoft.com/playwright:v$(shell sed -n 's/.*"@playwright\/test": "\([^"]*\)".*/\1/p' e2e/package.json)-noble
PLAYWRIGHT := docker run --rm --network host --ipc=host --user $$(id -u):$$(id -g) -e HOME=/tmp -v $(CURDIR)/e2e:/e2e -w /e2e
PO4A_IMAGE := debian:trixie-slim

.PHONY: dev prepare test docs clean

.env:
	cp .env.example .env

dev: .env web/node_modules internal/webui/dist/index.html
	@git config core.hooksPath .githooks
	go build -ldflags="-X main.version=$(VERSION)" -o bin/standblatt-dev ./cmd/standblatt
	@$(COMPOSE) up -d --wait db mail
	@set -a; . ./.env; set +a; ./bin/standblatt-dev & pid=$$!; \
		trap 'kill $$pid 2>/dev/null; $(COMPOSE) stop db mail' EXIT INT TERM; \
		for i in $$(seq 30); do curl -sf -o /dev/null http://localhost:8080/healthz && break; sleep 0.5; done; \
		$(COMPOSE) exec -T db psql -q -U standblatt -d standblatt < docker/dev-seed.sql >/dev/null; \
		echo "App: http://localhost:5173  login dev / devpassword, admin panel http://localhost:5173/admin (token in .env), mails http://localhost:8025"; \
		echo "Ctrl-C stops everything incl. Postgres"; \
		cd web && npm run dev

ERR_KEYWORDS := $(addprefix --keyword=,unauthorized forbidden missing conflict gone invalid tooMany)

# Everything to regenerate and check before a commit (the pre-commit hook in .githooks runs it): OpenAPI spec and
# TypeScript types, UI and document strings into the PO catalogs, gofmt, go vet, doc comments, svelte-check.
prepare: web/node_modules internal/webui/dist/index.html
	go run ./cmd/standblatt openapi > web/openapi.json
	cd web && npm run api
	cd web && npm run i18n
	xgettext -L Go --keyword=T:2 --keyword=Tf:2 --keyword=Problemf --keyword=notFound:2 $(ERR_KEYWORDS) --from-code=UTF-8 --add-location=file --no-wrap --sort-by-file \
		--package-name=standblatt -o internal/i18n/locales/messages.pot.go $$(find internal cmd -name '*.go' ! -name '*_test.go' | sort)
	sed -i.bak 's/charset=CHARSET/charset=UTF-8/' internal/i18n/locales/messages.pot.go && rm internal/i18n/locales/messages.pot.go.bak
	python3 scripts/tmplpot.py internal/i18n/locales/messages.pot.tmpl $$(find internal/web/templates -name '*.html' | sort)
	msgcat --no-wrap --sort-by-file --use-first -o internal/i18n/locales/messages.pot.new internal/i18n/locales/messages.pot.go internal/i18n/locales/messages.pot.tmpl
	rm internal/i18n/locales/messages.pot.go internal/i18n/locales/messages.pot.tmpl
	if diff -q -I POT-Creation-Date internal/i18n/locales/messages.pot.new internal/i18n/locales/messages.pot >/dev/null; \
		then rm internal/i18n/locales/messages.pot.new; else mv internal/i18n/locales/messages.pot.new internal/i18n/locales/messages.pot; fi
	for l in de fr it; do msgmerge --quiet --update --backup=none --no-wrap --no-fuzzy-matching internal/i18n/locales/$$l.po internal/i18n/locales/messages.pot; done
	@test -z "$$(gofmt -l cmd internal)" || { gofmt -l cmd internal; echo "gofmt: files above need formatting"; exit 1; }
	go vet ./...
	python3 scripts/funcdoc.py
	cd web && npm run check
	@git status --short

# Go unit and integration tests against the compose PostgreSQL, then the Playwright e2e specs against the real image.
test: internal/webui/dist/index.html
	@running=$$($(COMPOSE) ps -q --status running db); \
		$(COMPOSE) up -d --wait db; \
		TEST_DATABASE_URL=postgres://standblatt:standblatt@localhost:5433/standblatt go test ./...; status=$$?; \
		[ -n "$$running" ] || $(COMPOSE) stop db >/dev/null 2>&1; exit $$status
	$(COMPOSE) --profile e2e up --build --detach --force-recreate e2e
	rm -rf e2e/test-results
	@for i in $$(seq 30); do curl -sf -o /dev/null http://localhost:3002/healthz && break; sleep 1; done
	$(PLAYWRIGHT) $(PLAYWRIGHT_IMAGE) sh -c 'npm ci --no-audit --no-fund && npx playwright test --project=e2e'; \
		status=$$?; $(COMPOSE) --profile e2e rm --stop --force e2e e2e-db e2e-mail; exit $$status

web/node_modules: web/package-lock.json
	cd web && npm ci
	@touch web/node_modules

internal/webui/dist/index.html:
	mkdir -p internal/webui/dist && touch $@

# po4a and python3 in Debian; runs as root for apt, so the files are handed back to the user afterwards.
HELP_L10N = docker run --rm -t -v $(CURDIR)/web:/web -w /web/help $(PO4A_IMAGE) sh -c \
	'export DEBIAN_FRONTEND=noninteractive; apt-get update -qq && apt-get install -y -qq --no-install-recommends po4a python3 >/dev/null && python3 l10n.py $(1); \
	status=$$?; chown -R $$(stat -c %u:%g .) . ../src/help; exit $$status'

# Help: updates web/help/locales (POT and PO files) after the English pages changed and writes the translated pages
# (at least 80 % translated), then regenerates the screenshots from e2e/docs/screenshots.spec.js.
docs:
	$(call HELP_L10N,--update)
	$(COMPOSE) --profile e2e up --build --detach --force-recreate e2e
	@for i in $$(seq 30); do curl -sf -o /dev/null http://localhost:3002/healthz && break; sleep 1; done
	rm -rf web/static/help/screenshots && mkdir -p web/static/help/screenshots
	$(PLAYWRIGHT) -v $(CURDIR)/web/static/help/screenshots:/shots -e SHOTS=/shots $(PLAYWRIGHT_IMAGE) \
		sh -c 'npm ci --no-audit --no-fund && npx playwright test --project=docs'; \
		status=$$?; $(COMPOSE) --profile e2e rm --stop --force e2e e2e-db e2e-mail; exit $$status

# Stops the compose containers and removes build output and installed packages; keeps .env, database_seeds and the
# database volume (docker compose -f docker/compose.yaml down -v drops that too).
clean:
	-$(COMPOSE) --profile e2e down --remove-orphans
	rm -rf bin internal/webui/dist tmp web/node_modules web/.svelte-kit web/build web/src/locales/*.ts e2e/node_modules e2e/test-results
