package todoListRoute

import (
	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	appServices "github.com/Karik-ribasu/golang-todo-list-api/app/services"
	"github.com/labstack/echo/v4"
)

type controller struct {
	listItemService appServices.ListItemAppService
}

func newListItemController(listItemService appServices.ListItemAppService) *controller {
	return &controller{
		listItemService: listItemService,
	}
}

func (c controller) getListItems(ctx echo.Context) error {
	userUUID := ctx.Param("user-uuid")
	reqData := dto.GetListItemsRequest{
		UserUUID: userUUID,
	}

	resp, err := c.listItemService.GetListItems(reqData)
	if err != nil {
		return ctx.JSON(err.StatusCode, err.Message)
	}

	return ctx.JSON(200, resp)
}

func (c controller) postListItem(ctx echo.Context) error {
	userUUID := ctx.Param("user-uuid")

	reqData := dto.CreateListItemRequest{}
	err := ctx.Bind(&reqData)
	if err != nil {
		return ctx.JSON(422, `{"error": "Validation error}`)
	}
	reqData.UserUUID = userUUID

	resp, errHttp := c.listItemService.CreateListItem(reqData)
	if errHttp != nil {
		return ctx.JSON(errHttp.StatusCode, errHttp.Message)
	}

	return ctx.JSON(200, resp)
}

func (c controller) putListItem(ctx echo.Context) error {
	userUUID := ctx.Param("user-uuid")
	listItemUUID := ctx.Param("list-item-uuid")

	reqData := dto.UpdateListItemRequest{
		UserUUID:     userUUID,
		ListItemUUID: listItemUUID,
	}

	err := ctx.Bind(&reqData)
	if err != nil {
		return ctx.JSON(422, `{"error": "Validation error}`)
	}

	resp, errHttp := c.listItemService.UpdateListItem(reqData)
	if errHttp != nil {
		return ctx.JSON(errHttp.StatusCode, errHttp.Message)
	}

	return ctx.JSON(200, resp)
}

func (c controller) deleteListItem(ctx echo.Context) error {
	userUUID := ctx.Param("user-uuid")
	listItemUUID := ctx.Param("list-item-uuid")

	reqData := dto.DeleteListItemRequest{
		UserUUID:     userUUID,
		ListItemUUID: listItemUUID,
	}

	errHttp := c.listItemService.DeleteListItem(reqData)
	if errHttp != nil {
		return ctx.JSON(errHttp.StatusCode, errHttp.Message)
	}
	return ctx.NoContent(204)
}
