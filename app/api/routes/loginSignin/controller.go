package signInRoute

import (
	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	appServices "github.com/Karik-ribasu/golang-todo-list-api/app/services"
	"github.com/labstack/echo/v4"
)

type controller struct {
	loginSiginAppService appServices.LoginSiginAppService
}

func newLoginSiginController(loginSiginAppService appServices.LoginSiginAppService) *controller {
	return &controller{
		loginSiginAppService: loginSiginAppService,
	}
}

func (c controller) handleUserSignIn(ctx echo.Context) error {
	reqData := dto.SigninRequest{}
	err := ctx.Bind(&reqData)
	if err != nil {
		return ctx.JSON(422, `{"error": "Validation error}`)
	}

	errHttp := c.loginSiginAppService.SiginUser(reqData)
	if errHttp != nil {
		return ctx.JSON(errHttp.StatusCode, errHttp.Message)
	}

	return ctx.NoContent(204)
}

func (c controller) handleUserLogin(ctx echo.Context) error {
	reqData := dto.LoginRequest{}
	err := ctx.Bind(&reqData)
	if err != nil {
		return ctx.JSON(422, `{"error": "Validation error}`)
	}

	resp, errHttp := c.loginSiginAppService.LoginUser(reqData)
	if errHttp != nil {
		return ctx.JSON(errHttp.StatusCode, errHttp.Message)
	}

	return ctx.JSON(200, resp)
}
