package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func moduleRoot(t *testing.T) string {
	t.Helper()
	wd, _ := os.Getwd()
	dir := wd
	for i := 0; i < 15; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("root")
	return ""
}

func writeRunConfig(t *testing.T, dir, pemAbs string) {
	t.Helper()
	content := "[db]\nuser = \"a\"\npasswd = \"b\"\naddr = \"c\"\nport = \"3306\"\nname = \"d\"\n\n[app]\ncertificate_key_path = \"" + pemAbs + "\"\n"
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRun_ConfigMissing(t *testing.T) {
	dir := t.TempDir()
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := run(); err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_DBInitFails(t *testing.T) {
	root := moduleRoot(t)
	dir := t.TempDir()
	writeRunConfig(t, dir, filepath.Join(root, "testdata", "dev_rsa_private.pem"))
	wd, _ := os.Getwd()
	origDB := dbInit
	origStart := startEcho
	t.Cleanup(func() {
		_ = os.Chdir(wd)
		dbInit = origDB
		startEcho = origStart
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	dbInit = func(config.Config) (data.DbManager, error) { return nil, errors.New("db") }
	startEcho = func(e *echo.Echo) error { return nil }
	if err := run(); err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_StartEcho(t *testing.T) {
	root := moduleRoot(t)
	dir := t.TempDir()
	writeRunConfig(t, dir, filepath.Join(root, "testdata", "dev_rsa_private.pem"))
	wd, _ := os.Getwd()
	origDB := dbInit
	origStart := startEcho
	t.Cleanup(func() {
		_ = os.Chdir(wd)
		dbInit = origDB
		startEcho = origStart
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	dbInit = func(config.Config) (data.DbManager, error) { return data.NewManagerFromDB(db), nil }
	startEcho = func(e *echo.Echo) error { return nil }
	if err := run(); err != nil {
		t.Fatal(err)
	}
}

func TestRun_StartEchoError(t *testing.T) {
	root := moduleRoot(t)
	dir := t.TempDir()
	writeRunConfig(t, dir, filepath.Join(root, "testdata", "dev_rsa_private.pem"))
	wd, _ := os.Getwd()
	origDB := dbInit
	origStart := startEcho
	t.Cleanup(func() {
		_ = os.Chdir(wd)
		dbInit = origDB
		startEcho = origStart
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	dbInit = func(config.Config) (data.DbManager, error) { return data.NewManagerFromDB(db), nil }
	startEcho = func(*echo.Echo) error { return errors.New("listen") }
	if err := run(); err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_InvalidRSA(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte("[db]\nuser=a\npasswd=b\naddr=c\nport=1\nname=d\n\n[app]\ncertificate_key=\"not-pem\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := run(); err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_LoadCertificateError(t *testing.T) {
	dir := t.TempDir()
	writeRunConfig(t, dir, filepath.Join(dir, "missing.pem"))
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := run(); err == nil {
		t.Fatal("expected error")
	}
}

func TestMainInvokesMainRun(t *testing.T) {
	prev := mainRun
	defer func() { mainRun = prev }()
	mainRun = func() error { return nil }
	main()
}

func TestMain_LogFatalfOnError(t *testing.T) {
	var saw bool
	prevL := logFatalf
	prevR := mainRun
	defer func() {
		logFatalf = prevL
		mainRun = prevR
	}()
	logFatalf = func(string, ...any) { saw = true }
	mainRun = func() error { return errors.New("e") }
	main()
	if !saw {
		t.Fatal("expected logFatalf")
	}
}

func TestDefaultLogFatalf(t *testing.T) {
	prev := osExit
	defer func() { osExit = prev }()
	osExit = func(int) {}
	defaultLogFatalf("err: %v", errors.New("x"))
}

func TestDefaultStartEcho_ListenShutdown(t *testing.T) {
	t.Setenv("TODO_LISTEN_ADDR", "127.0.0.1:0")
	e := echo.New()
	go func() {
		time.Sleep(40 * time.Millisecond)
		_ = e.Shutdown(context.Background())
	}()
	_ = defaultStartEcho(e)
}

func TestListenAddr_FromEnv(t *testing.T) {
	t.Setenv("TODO_LISTEN_ADDR", ":4001")
	if listenAddr() != ":4001" {
		t.Fatal(listenAddr())
	}
}

func TestListenAddr_Default(t *testing.T) {
	t.Setenv("TODO_LISTEN_ADDR", "")
	if listenAddr() != ":5000" {
		t.Fatal(listenAddr())
	}
}

func TestRun_RSAInvalidPEM(t *testing.T) {
	dir := t.TempDir()
	bad := []byte("-----BEGIN PRIVATE KEY-----\nYWJ=\n-----END PRIVATE KEY-----\n")
	pemPath := filepath.Join(dir, "bad.pem")
	if err := os.WriteFile(pemPath, bad, 0o600); err != nil {
		t.Fatal(err)
	}
	writeRunConfig(t, dir, pemPath)
	wd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	if err := run(); err == nil {
		t.Fatal("expected parse rsa error")
	}
}

func TestInitializeServer_JWTBlocksTodoWithoutToken(t *testing.T) {
	root := moduleRoot(t)
	b, err := os.ReadFile(filepath.Join(root, "testdata", "dev_rsa_private.pem"))
	if err != nil {
		t.Fatal(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{App: config.App{PrivateKey: key}}
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	e := InitializeServer(cfg, data.NewManagerFromDB(db))
	req := httptest.NewRequest(http.MethodGet, "/user/u/todo-list", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatalf("unexpected %d", rec.Code)
	}
}
