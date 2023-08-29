package todoListRoute

import (
	appServices "github.com/Karik-ribasu/golang-todo-list-api/app/services"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, listItemService appServices.ListItemAppService) {

	controller := newListItemController(listItemService)

	e.GET("/user/:user-uuid/todo-list", controller.getListItems)
	e.POST("/user/:user-uuid/todo-list", controller.postListItem)
	e.PUT("/user/:user-uuid/todo-list/:list-item-uuid", controller.putListItem)
	e.DELETE("/user/:user-uuid/todo-list/:list-item-uuid", controller.deleteListItem)
}
