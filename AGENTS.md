# AGENTS.md

## Cursor Cloud specific instructions

### Overview

This is a **Go REST API** for a todo list application using **Echo v4**, **MySQL**, and **JWT authentication**. It runs on port **5000**.

### Prerequisites (one-time setup per VM)

MySQL 8.0 must be installed and running. The app requires:

1. A running MySQL instance on `localhost:3306`
2. A `config.toml` file at the project root (gitignored) with DB credentials and RSA key path
3. A dev RSA private key is already committed at `testdata/dev_rsa_private.pem` (PKCS#8 format)

### Starting MySQL

```bash
sudo mkdir -p /var/run/mysqld && sudo chown mysql:mysql /var/run/mysqld
sudo mysqld --user=mysql --datadir=/var/lib/mysql &
sleep 3
```

### Database setup

The database name is `todo` (not `todolist`). Use the schema from `scripts/init.sql`:

```bash
sudo mysql -u root -e "
CREATE DATABASE IF NOT EXISTS todo;
USE todo;
SOURCE /workspace/scripts/init.sql;
CREATE USER IF NOT EXISTS 'todo'@'localhost' IDENTIFIED BY 'todo';
GRANT ALL PRIVILEGES ON todo.* TO 'todo'@'localhost';
FLUSH PRIVILEGES;
"
```

### config.toml

Copy from the example: `cp config.example.toml config.toml`. The example already points to the dev RSA key:

```toml
[db]
user = "todo"
passwd = "todo"
addr = "127.0.0.1"
port = "3306"
name = "todo"

[app]
certificate_key_path = "./testdata/dev_rsa_private.pem"
```

### Standard commands

| Action | Command |
|--------|---------|
| Build | `go build -o todoapi .` |
| Run | `go run .` |
| Lint | `go vet ./...` |
| Test (unit) | `go test ./...` |
| Test (with coverage) | `make test-unit` |
| Integration tests | `make test-integration` (needs MySQL + `config.toml`; CI uses `config.ci.toml`) |
| Dev in Docker (API + DB) | `docker compose up --build` |
| Dev shell (Go + source + DB) | `make docker-dev-up` then `docker compose -f docker-compose.yml -f docker-compose.dev.yml exec dev bash` |
| Dependencies | `go mod download` |

### Notes

- The JWT middleware in `initServer.go` has a Skipper that exempts `/sign-in` and `/log-in`.
- In `domain/data/db.go`, the DSN uses `Net: "tcp"` and `Addr` as `cfg.Db.Addr + ":" + cfg.Db.Port` so the driver emits `tcp(host:port)/dbname` (required for non-default ports and matches Docker service names such as `mysql:3306`).
