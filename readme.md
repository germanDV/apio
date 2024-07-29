# APIO

Template for JSON-HTTP APIs.

Features:
- [X] Auth with JWT, including RBAC for two simple roles: _user_ and _admin_.
- [ ] IP rate limiting using Redis.
- [ ] Data persistence using Postgres, including migrations.
- [X] _.env_ file support.
- [ ] OpenAPI 3 docs and Swagger UI.
    * https://github.com/parvez3019/go-swagger3 -> generates spec but doesn't have UI
    * https://github.com/a-h/rest -> specs + UI
- [ ] Tracing with OpenTelemetry (Grafana + Tempo).
- [ ] Metrics (Prometheus).
- [ ] Logging (Loki).
- [ ] Github Actions for PRs.
- [X] Air for hot-reloading during development.
- [ ] Docker to easily run everything locally.
- [X] A Makefile for convenience.

## Design Decisions

- Simple to understand, avoid too much abstraction and dependency on other libraries that abstract things away.
- Minimal dependencies, so that it is easier to keep up-to-date with library upgrades, there's less "external" code and potentially a smaller surface for security vulnerabilities.
- Follow well-known Go idioms to make it familiar to most Go developers.
- Implement only things that I have used in several projects and found myself copy-pasting. Do not try to build a fully customizable framework to fit every need.

## Usage

In order not to pollute the template with many things that would need to be removed, a single entity has been created with basic CRUD functionality.

## Run It Locally

Start docker containers for all dependencies:
```sh
# TODO
```

Run API server in _DEV_ mode:
```sh
# TODO
```
