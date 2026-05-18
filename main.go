package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func listenAddr() string {
	addr := os.Getenv("TODO_LISTEN_ADDR")
	if addr == "" {
		return ":5000"
	}
	return addr
}

func defaultStartEcho(e *echo.Echo) error {
	return e.Start(listenAddr())
}

var (
	dbInit    = data.InitDB
	startEcho = defaultStartEcho
	mainRun   = run
	osExit    = os.Exit
)

func defaultLogFatalf(format string, v ...any) {
	log.Printf(format, v...)
	osExit(1)
}

var logFatalf = defaultLogFatalf

func main() {
	if err := mainRun(); err != nil {
		logFatalf("%v", err)
	}
}

func run() error {
	cfg, err := config.ReadConfig()
	if err != nil {
		return err
	}

	pemBytes, err := config.LoadCertificatePEM(cfg.App)
	if err != nil {
		return err
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return fmt.Errorf("parse rsa key: %w", err)
	}
	cfg.App.PrivateKey = rsaKey

	db, err := dbInit(cfg)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}

	e := InitializeServer(cfg, db)
	return startEcho(e)
}
