# AGENTS.md

## Overview

This is a **Go REST API** for a todo list application using **Echo v4**, **MySQL 8**, and **JWT (RS256)**. The HTTP server listens on `:5000` by default, or on `TODO_LISTEN_ADDR` when set.

Authoritative procedure for agents lives under **`.cursor/agents/`** and **`.cursor/skills/`**.

## Containerized development

### Dev Container (recommended)

Open the repository in VS Code or Cursor and use **“Dev Container: Reopen in Container”**. The stack is defined in `.devcontainer/` and includes the Go toolchain plus MySQL with `scripts/init.sql` applied on first start.

After the container is created, `config.toml` is generated from `config.docker.toml` (`postCreateCommand`). Run the API with:

```bash
go run .
```

### Docker Compose (host or CI-like)

| Command | Services | API URL |
|---------|----------|---------|
| `docker compose up -d` | MySQL + distroless API (image build) | http://localhost:5000 |
| `docker compose --profile dev up -d` | MySQL + **api-dev** (`go run` in a Go image) | http://localhost:5001 |

Database name: **`todo`**. Credentials match `config.docker.toml` / `config.ci.toml`.

## Configuration

Copy `config.example.toml` to **`config.toml`** (gitignored) or use **`config.ci.toml`** / **`config.docker.toml`** as templates.

- Either **`app.certificate_key_path`** (PEM file, PKCS#1 RSA private key) or **`app.certificate_key`** (inline PEM) must be set.
- For local unit tests, the repo ships **`testdata/dev_rsa_private.pem`** (development only).

## Database schema

Apply **`scripts/init.sql`** to the `todo` database (Compose mounts it for MySQL; CI runs `mysql ... < scripts/init.sql`).

## Commands

| Action | Command |
|--------|---------|
| Lint | `go vet ./...` or `make vet` |
| Unit tests + 100% module coverage | `make test-unit` |
| Integration tests | `make test-integration` (needs DB + `config.toml`) |
| Merge coverage profiles | `make merge-coverage` |
| Build | `go build -trimpath -ldflags="-s -w" -o bin/todo-api .` |
| Dependencies | `go mod download` |

## HTTP behavior (sanity)

- **JWT** middleware skips **`/sign-in`** and **`/log-in`**; other routes require `Authorization: Bearer <token>`.
- **`POST /sign-in`** registers a user; **`POST /log-in`** returns a JWT.

## CI workflows (GitHub Actions)

| Workflow | Purpose |
|----------|---------|
| `unit-tests.yml` | `go vet`, unit tests, **100%** statements (`-coverpkg=./...`) |
| `integration-tests.yml` | MySQL service, schema, unit + `-tags=integration` tests, merged coverage, **100%** total |
| `build.yml` | `go vet`, static Linux binary artifact |
| `release.yml` | On `v*` tags: build + **GitHub Release** with binaries |

## DB address format

In `domain/data/db.go`, the DSN address is built as `cfg.Db.Addr + ":" + cfg.Db.Port` for the MySQL driver (TCP).
