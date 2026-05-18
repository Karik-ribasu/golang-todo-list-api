# Skill: testing, coverage, and CI for this repo

## Unit tests

1. Run from repository root: `go test ./... -count=1 -coverprofile=coverage-unit.out -covermode=atomic -coverpkg=./...`
2. Inspect total: `go tool cover -func=coverage-unit.out | tail -1`
3. CI job `unit-tests` must report `100.0%` on that profile.

## Integration tests

1. Build tag: `integration` (`//go:build integration` in `main_integration_test.go`).
2. Requires MySQL with schema applied and `config.toml` (or copy `config.ci.toml` → `config.toml` for local TCP to `127.0.0.1:3306`).
3. Command: `go test -tags=integration -count=1 -timeout=5m -coverprofile=coverage-integration.out -covermode=atomic -coverpkg=./... .`
4. Only the root `main` package carries integration tests; keep `-coverpkg=./...` so the merged report spans the module.

## Merged coverage (integration pipeline)

1. Produce `coverage-unit.out` and `coverage-integration.out` as above.
2. Merge: `python3 scripts/merge_cover_profiles.py coverage-merged.out coverage-unit.out coverage-integration.out`
3. Gate: `go tool cover -func=coverage-merged.out | tail -1` must show `100.0%`.

## Containerized development

1. API + MySQL: `docker compose up --build` (see root `docker-compose.yml`).
2. Go dev shell with compose: `docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d` then `exec` into `dev` (see `Makefile` target `docker-dev-up`).
3. VS Code / Cursor Dev Containers: open repo with `.devcontainer/devcontainer.json`.
