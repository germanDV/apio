## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## dev: run with hot-reloading
.PHONY: dev
dev:
	air .

## test: run tests
.PHONY: test
test:
	go test ./...

## docker/up: start PostgreSQL + Redis docker containers
.PHONY: docker/up
docker/up:
	@echo 'Starting docker-compose'
	docker compose up -d

## docker/stop: stop docker containers
.PHONY: docker/stop
docker/stop:
	@echo 'Stopping docker-compose'
	docker compose stop

## docker/down: tear down docker containers
.PHONY: docker/down
docker/down: confirm
	@echo 'Stopping docker-compose'
	docker compose down

## db/migrate/up: run database migrations
.PHONY: db/migrate/up
db/migrate/up:
	@echo 'Running migrations...'
	@go run ./cmd/migrate -action up

## db/migrate/down: rollback latest database migration
.PHONY: db/migrate/down
db/migrate/down: confirm
	@echo 'Rolling back latest migration..'
	@go run ./cmd/migrate -action down

## db/cli: connect to local database using pgcli
.PHONY: db/cli
db/cli:
	@echo 'Connecting to database...'
	pgcli -h localhost -p 5432 -U postgres -d apio

## deps/upgrade/all: upgrade all dependencies
.PHONY: deps/upgrade/all
deps/upgrade/all:
	@echo 'Upgrading dependencies to latest versions...'
	go get -t -u ./...

## deps/upgrade/patch: upgrade dependencies to latest patch version
.PHONY: deps/upgrade/patch
deps/upgrade:
	@echo 'Upgrading dependencies to latest patch versions...'
	go get -t -u=patch ./...

## vuln: check for vulnerabilities
.PHONY: vuln
vuln:
	govulncheck ./...

## vet: vet code
.PHONY: vet
vet:
	@echo 'Vetting code...'
	go vet ./...

## lint: lint code
.PHONY: lint
lint:
	@echo 'Linting code with golangci-lint...'
	golangci-lint run --disable-all --enable errcheck,gosimple,ineffassign,unused,gocritic,misspell,stylecheck ./...

## sec: check for security issues
.PHONY: sec
sec:
	@echo 'Running security checks with gosec...'
	gosec ./...

## scripts/token u=$1 r=$2: generate an auth token ($ make scripts/token u=<user_id> r=<role>)
.PHONY: scripts/token
scripts/token:
	go run ./cmd/token "$u" "$r"
