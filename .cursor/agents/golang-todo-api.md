# Agent role: golang-todo-list-api

Authoritative context for automated agents working in this repository.

## Stack and runtime

- Go REST API with Echo v4, MySQL 8, JWT (RSA RS256).
- Default listen address `:5000` unless `TODO_LISTEN_ADDR` is set.
- Configuration is TOML at `./config.toml` (gitignored); use `config.example.toml` as a template.

## Database

- Database name: `todo`. Apply `scripts/init.sql` to a fresh database.
- MySQL DSN is built in `domain/data/db.go` with explicit `Net: "tcp"` so host/port from config are honored in the driver DSN.

## Commands

| Goal | Command |
|------|---------|
| Build | `go build -o todoapi .` |
| Run | `go run .` |
| Lint | `go vet ./...` |
| Unit tests + coverage | `make test-unit` |
| Integration tests + coverage | `make test-integration` (requires MySQL and valid `config.toml`; CI uses `config.ci.toml`) |

## Coverage and CI

- Statement coverage uses `-coverpkg=./...` across the module.
- The integration workflow merges unit and integration profiles (`scripts/merge_cover_profiles.py`) and requires 100% total after merge.

## Git and branches

- Implement on the assigned feature branch; push only to that branch unless explicitly told otherwise.
