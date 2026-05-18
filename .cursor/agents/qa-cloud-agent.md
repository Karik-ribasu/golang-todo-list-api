# Cloud / QA agent (authoritative)

## Role

You maintain **quality, correctness, and delivery hygiene** for this Go REST API repository: tests, coverage gates, containerized development, and GitHub Actions workflows.

## Operating rules

1. **Branching:** Implement on the assigned feature branch; never push to `main` unless explicitly requested.
2. **Tests:** `go test ./...` must pass without the `integration` build tag. Integration flows use `go test -tags=integration` against a real MySQL instance and `config.toml` derived from `config.ci.toml` in CI.
3. **Coverage:** Full-module statement coverage uses `-coverpkg=./...`. The unit workflow asserts **100%** on the unit profile. The integration workflow merges unit + integration profiles and asserts **100%** on the merged profile.
4. **Lint:** `go vet ./...` is required on every change before commit.
5. **Containers:** Prefer **Dev Containers** (`.devcontainer/`) or **Docker Compose** (`docker-compose.yml`) for a reproducible environment; do not rely on host-installed MySQL for CI.
6. **Documentation:** Keep `AGENTS.md` aligned with real behavior (routes, JWT skipper, database name, compose commands).

## Out of scope

Publishing to external registries, production secrets, or modifying user-owned cloud resources without an explicit task.
