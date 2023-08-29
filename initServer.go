package main

import (
	loginSigninRoute "github.com/Karik-ribasu/golang-todo-list-api/app/api/routes/loginSignin"
	todoListRoute "github.com/Karik-ribasu/golang-todo-list-api/app/api/routes/todoList"
	appServices "github.com/Karik-ribasu/golang-todo-list-api/app/services"
	data "github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	domainServices "github.com/Karik-ribasu/golang-todo-list-api/domain/services"
	config "github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeServer(cfg config.Config, db data.DbManager) (e *echo.Echo) {
	// Echo instance
	e = echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    cfg.App.PrivateKey,
		SigningMethod: echojwt.AlgorithmHS256,
		TokenLookup:   "header:Authorization:Bearer",
	}))

	// init domain services
	domainServices := domainServices.NewDomainSVC(db)

	// init app services
	appServices := appServices.NewAppService(cfg, domainServices)

	// init routes
	todoListRoute.Init(e, appServices.TodoListAppService())
	loginSigninRoute.Init(e, appServices.LoginSiginAppService())

	return e
}
