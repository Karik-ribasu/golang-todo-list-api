package main

import (
	"log"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	jwt "github.com/golang-jwt/jwt/v4"
)

func main() {

	cfg, err := config.ReadConfig()

	// init DB
	db, err := data.InitDB(cfg)
	if err != nil {
		log.Fatal("error connecting to database: ", err.Error())
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.App.CertificateKey))
	if err != nil {
		log.Fatal("error getting rsaKey: ", err.Error())
	}
	cfg.App.PrivateKey = rsaKey

	e := InitializeServer(cfg, db)

	// Start server
	e.Logger.Fatal(e.Start(":5000"))
}
