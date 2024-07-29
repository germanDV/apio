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

## scripts/token u=$1 r=$2: generate an auth token ($ make scripts/token u=<user_id> r=<role>)
.PHONY: scripts/token
scripts/token:
	go run ./cmd/token "$u" "$r"
