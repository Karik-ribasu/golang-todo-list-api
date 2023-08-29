package signInRoute

import (
	appServices "github.com/Karik-ribasu/golang-todo-list-api/app/services"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, loginSiginAppService appServices.LoginSiginAppService) {

	controller := newLoginSiginController(loginSiginAppService)

	e.POST("/sign-in", controller.handleUserSignIn)
	e.POST("/log-in", controller.handleUserSignIn)
}
