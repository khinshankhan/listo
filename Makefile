# Set .SHELLFLAGS to use "bash strict mode". This will fail on any error, fail
# on an error in a pipeline (usually it just returns the value of the last
# command in a pipeline), and will fail if using any undefined variables.
SHELL := bash
# Set .ONESHELL config. Runs the whole make recipe in one shell session.
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
# Both of the above are GNU Make 3.82+. macOS ships 3.81 and ignores them,
# so a multi-command recipe there gets no strict mode and needs its own
# `set -e` to keep a failure from passing as success.
# Set .DELETE_ON_ERROR. Deletes any files generated on error.
.DELETE_ON_ERROR:
# Set makeflags to --warn-undefined-variables (probably an error) and to avoid
# the built-in rules (this removes a lot of magic that is related to yacc and
# other tools).
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# cgo free on every platform, so an accidental cgo dependency fails the build.
# `test` turns it back on for itself since the race detector needs it.
export CGO_ENABLED := 0

# VARIABLES
#
# A service is any ./cmd/*/main.go.
GOSERVICES = $(sort $(notdir $(realpath $(dir $(wildcard ./cmd/*/main.go)))))

.PHONY: $(GOSERVICES)

.PHONY: all
all: check build

# web and templ output feed the go build, so they go first
.PHONY: build
build: generate $(GOSERVICES)
	go build ./...

# creates a TARGET per go service
$(GOSERVICES): % : ./cmd/%/main.go generate
	go build -o bin/$@ ./cmd/$@

# list available go services
.PHONY: services
services:
	@echo $(GOSERVICES)

.PHONY: generate
generate: web templ

# WEB
#
# node only runs at build time, the bundle is embedded in the binary
.PHONY: web
web: web/node_modules
	cd web
	pnpm exec tsc --noEmit
	pnpm exec vite build

web/node_modules: web/package.json web/pnpm-lock.yaml
	cd web
	pnpm install --frozen-lockfile
	touch node_modules

.PHONY: templ
templ:
	go tool templ generate

# build everything, then serve
.PHONY: serve
serve: listo
	./bin/listo serve

# serve, rebuilding on every save and reloading the browser: open DEV_PORT,
# air's proxy, which forwards to listo on the port after it
DEV_PORT ?= 8080
.PHONY: dev
dev: db-up web/node_modules
	app=$$(($(DEV_PORT) + 1))
	LISTO_ADDR=127.0.0.1:$$app go tool air --proxy.proxy_port $(DEV_PORT) --proxy.app_port $$app

# DATABASE
#
# Postgres 18 from the flake (`direnv allow` or `nix develop`), as a cluster
# under var/pg on a unix socket only, so it can't collide with anything.
# Not secrets; .envrc sets the same.
export PGDATA := $(CURDIR)/var/pg
export PGHOST := $(PGDATA)
export PGDATABASE := listo

# fail early with a hint, rather than "pg_ctl: command not found"
.PHONY: need-pg
need-pg:
	@command -v pg_ctl >/dev/null || { echo "no postgres on PATH: run direnv allow, or nix develop"; exit 1; }

$(PGDATA)/PG_VERSION: | need-pg
	initdb --auth=trust --encoding=UTF8 --locale=C.UTF-8 --no-instructions >/dev/null

# start the cluster if needed and make sure the database exists.
# NOTE: -X skips ~/.psqlrc, whose output (eg \timing) would break the check
.PHONY: db-up
db-up: $(PGDATA)/PG_VERSION
	@pg_ctl status >/dev/null || pg_ctl start --wait --log=$(PGDATA)/server.log \
		-o "-k $(PGDATA) -c listen_addresses=''"
	@psql -X -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = '$(PGDATABASE)'" | grep -qx 1 \
		|| createdb $(PGDATABASE)

# rebuild the database from schema/ and seed/; no migrations until there's
# data worth keeping
.PHONY: db-reset
db-reset: db-up
	dropdb --force $(PGDATABASE)
	createdb $(PGDATABASE)
	psql -X -q -v ON_ERROR_STOP=1 -f schema/schema.sql
	for f in seed/*.sql; do psql -X -q -v ON_ERROR_STOP=1 -f "$$f" >/dev/null; done

.PHONY: db-down
db-down: need-pg
	@pg_ctl status >/dev/null && pg_ctl stop --wait || echo "not running"

.PHONY: psql
psql: db-up
	@psql

# regenerate the typed query layer from schema/ + queries/
.PHONY: sqlc
sqlc: | need-pg
	sqlc generate

# fails instead of rewriting, so CI catches stale generated code
.PHONY: sqlc-check
sqlc-check: | need-pg
	sqlc diff

# extra flags and a package selector for the test targets, eg
#   make test GOTESTFLAGS=-v
#   make test GOTESTFLAGS='-run TestFoo -v'
#   make test PKG=./controller/...
GOTESTFLAGS ?=
PKG ?= ./...

# -race by default, which needs cgo, hence the local CGO_ENABLED override.
# data/ tests need the dev database loaded once (make db-reset)
.PHONY: test
test: generate db-up
	CGO_ENABLED=1 go test -race $(GOTESTFLAGS) $(PKG)

# no race detector, for a quick inner loop or a cgo-less environment
.PHONY: test-quick
test-quick: generate db-up
	go test $(GOTESTFLAGS) $(PKG)

.PHONY: check
check: fmt-check vet sqlc-check test

.PHONY: fmt
fmt:
	gofmt -w .
	go tool templ fmt views
	go mod tidy

# fails instead of rewriting, so CI catches unformatted code
.PHONY: fmt-check
fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "not gofmt'd:"; echo "$$unformatted"; exit 1; \
	fi

.PHONY: vet
vet: generate
	go vet ./...

.PHONY: clean
clean:
	rm -rf bin web/dist/bundle
	find views -name '*_templ.go' -delete
