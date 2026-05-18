# Skill: Go tests, coverage, and CI (authoritative)

## Unit tests and coverage

From the repository root:

```bash
go vet ./...
go test ./... -count=1 -coverprofile=coverage-unit.out -covermode=atomic -coverpkg=./...
go tool cover -func=coverage-unit.out | tail -1
```

The last line must report **100.0%** total statements for the `unit-tests` workflow gate.

## Integration tests and merged coverage

Requires MySQL 8 with schema from `scripts/init.sql`, database `todo`, and `config.toml` (CI uses `cp config.ci.toml config.toml`).

```bash
go test -tags=integration -count=1 -timeout=5m -coverprofile=coverage-integration.out -covermode=atomic -coverpkg=./... .
python3 scripts/merge_cover_profiles.py coverage-merged.out coverage-unit.out coverage-integration.out
go tool cover -func=coverage-merged.out | tail -1
```

The merged total must be **100.0%** for the `integration-tests` workflow gate.

## Local containers

- **VS Code / Cursor Dev Container:** open the repo folder and “Reopen in Container” — see `.devcontainer/devcontainer.json`.
- **Compose (distroless API):** `docker compose up -d` — API on port **5000**.
- **Compose (dev API with `go run`):** `docker compose --profile dev up -d` — API on host port **5001**.

## Pipelines (GitHub Actions)

| Workflow | Purpose |
|----------|---------|
| `unit-tests.yml` | `go vet`, unit tests, 100% coverage (unit profile) |
| `integration-tests.yml` | MySQL service, schema, unit + integration tests, merge coverage, 100% merged |
| `build.yml` | `go vet`, `go build` artifact |
| `release.yml` | Tag builds, GitHub Release + binaries |
