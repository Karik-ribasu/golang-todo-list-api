# AGENTS.md

## Cursor Cloud specific instructions

### Overview

This is a **Go REST API** for a todo list application using **Echo v4**, **MySQL**, and **JWT authentication**. It runs on port **5000**.

### Prerequisites (one-time setup per VM)

MySQL 8.0 must be installed and running. The app requires:

1. A running MySQL instance on `localhost:3306`
2. A `config.toml` file at the project root (gitignored) with DB credentials and an RSA private key
3. An RSA private key in PKCS#1 format (`-----BEGIN RSA PRIVATE KEY-----`)

### Starting MySQL

```bash
sudo mkdir -p /var/run/mysqld && sudo chown mysql:mysql /var/run/mysqld
sudo mysqld --user=mysql --datadir=/var/lib/mysql &
sleep 3
```

### Database setup

```bash
sudo mysql -u root -e "
CREATE DATABASE IF NOT EXISTS todolist;
USE todolist;
CREATE TABLE IF NOT EXISTS user (
    user_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_uuid VARCHAR(36) NOT NULL UNIQUE,
    nick_name VARCHAR(255) NOT NULL UNIQUE,
    password BLOB NOT NULL,
    active BOOLEAN DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS list_item (
    list_item_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    list_item_uuid VARCHAR(36) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (user_id) REFERENCES user(user_id)
);
CREATE USER IF NOT EXISTS 'todouser'@'localhost' IDENTIFIED BY 'todopass';
GRANT ALL PRIVILEGES ON todolist.* TO 'todouser'@'localhost';
FLUSH PRIVILEGES;
"
```

### Generating RSA key and config.toml

```bash
openssl genrsa -traditional -out cert.pem 2048
```

The `config.toml` should have this structure (both files are gitignored):

```toml
[db]
user = "todouser"
passwd = "todopass"
addr = "127.0.0.1"
port = "3306"
name = "todolist"

[app]
certificate_key = "<PEM key content with \\n for newlines>"
```

### Standard commands

| Action | Command |
|--------|---------|
| Build | `go build -o todoapi .` |
| Run | `go run .` |
| Lint | `go vet ./...` |
| Test | `go test ./...` |
| Dependencies | `go mod download` |

### Known code issues (not environment issues)

1. **Global JWT middleware**: The JWT middleware in `initServer.go` is applied to all routes including `/sign-in` and `/log-in`. These public routes return 401 because no `Skipper` is configured. This is a code bug, not an env issue.
2. **Login route mapping**: In `app/api/routes/loginSignin/route.go`, the `/log-in` route maps to `handleUserSignIn` instead of `handleUserLogin`.
3. **No test files**: The codebase has zero test files.

### DB address format gotcha

In `domain/data/db.go`, the MySQL address is built as `cfg.Db.Addr + "::" + cfg.Db.Port` (note the double colon). The Go MySQL driver interprets this as a TCP address when configured through `mysql.NewConfig()`, so using `127.0.0.1` as `addr` and `3306` as `port` works correctly.
