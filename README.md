# golang-todo-list-api

Todo list HTTP API in Go (Echo, MySQL, JWT).

## Configuration

Place `config.toml` at the repository root (and the certificate or key material your config references); those paths are gitignored. Start from `config.example.toml`.

Run the server with `go run .` (default listen `:5000`, overridable with `TODO_LISTEN_ADDR`).

## Containerized development

- **API and MySQL only:** `docker compose up --build` — uses `config.docker.toml` (MySQL host `mysql`) and the dev key mounted in compose.
- **Go toolchain + mounted source** (for `go test`, `go run`, etc. against the compose network): `make docker-dev-up`, then attach with  
  `docker compose -f docker-compose.yml -f docker-compose.dev.yml exec dev bash`.  
  Inside the container, copy or mount `config.docker.toml` as `config.toml` if needed.
- **VS Code / Cursor Dev Containers:** open the repo; use `.devcontainer/devcontainer.json` (same compose files).

## Tests and CI

- Unit coverage: `make test-unit` (`-coverpkg=./...`).
- Integration (real MySQL): `make test-integration` with `-tags=integration` (see `AGENTS.md` for database setup).
- CI runs separate workflows for unit tests, integration tests (merged coverage gate), build, and release on version tags.
