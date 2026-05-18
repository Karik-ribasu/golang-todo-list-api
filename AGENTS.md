# AGENTS.md

## Cursor Cloud specific instructions

### Overview

This is a **Go Todo List REST API** built with [Echo](https://echo.labstack.com/) v4, backed by MySQL, with JWT (RSA) authentication. Single service, no monorepo. Runs on port **5000**.

### Prerequisites (already installed in VM)

- **Go 1.20+** (system Go at `/usr/bin/go`)
- **MySQL 8.0** (started manually via `sudo mysqld --user=mysql &`)
- **staticcheck** installed at `~/go/bin/staticcheck`

### Starting MySQL

MySQL must be started manually each session:

```bash
sudo mkdir -p /var/run/mysqld && sudo chown mysql:mysql /var/run/mysqld
sudo mysqld --user=mysql &
sleep 3
sudo chmod 755 /var/run/mysqld/ && sudo chmod 777 /var/run/mysqld/mysqld.sock
```

Database `todolist` with tables `user` and `list_item` is pre-created. User: `todouser` / `todopass`.

### Configuration

The app reads `config.toml` from project root (gitignored). It also requires an RSA private key. Both are pre-created in the VM. If regenerating:

```bash
openssl genrsa -out cert.pem 2048
```

Then embed the PEM content (with `\n` escaped newlines) in `config.toml` under `[app] certificate_key`.

### Build / Lint / Test / Run

| Action | Command |
|--------|---------|
| Build | `go build -o todolist-api .` |
| Lint (vet) | `go vet ./...` |
| Lint (staticcheck) | `~/go/bin/staticcheck ./...` |
| Test | `go test ./...` (no test files exist in the codebase) |
| Run | `go run .` (listens on `:5000`) |

### Known code-level issues

1. **JWT middleware applied globally** (`initServer.go`): The `echojwt` middleware is applied via `e.Use()` before routes, including `/sign-in` and `/log-in`. These auth routes return 401 because no `Skipper` is configured. A fix would add a `Skipper` function to exclude auth routes.
2. **JWT algorithm mismatch** (`initServer.go` vs `infra/auth/jwt.go`): Tokens are signed with `RS256` but the middleware validates with `HS256`. The middleware `SigningKey` is an RSA private key but configured for HMAC.
3. **DB address double colon** (`domain/data/db.go:19`): The address is built as `addr + "::" + port` (double colon). This doesn't cause connection failures because `mysql.NewConfig()` leaves `Net` empty, so `FormatDSN()` omits the address entirely and the driver falls back to Unix socket.
4. **Login route handler mismatch** (`app/api/routes/loginSignin/route.go:13`): Both `/sign-in` and `/log-in` are mapped to `handleUserSignIn`; `/log-in` should map to `handleUserLogin`.
