package todoListRoute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/errors"
	"github.com/labstack/echo/v4"
)

type fakeListSvc struct {
	g  dto.GetListItemsResponse
	ge *errors.HttpError
	cr dto.CreateListItemResponse
	ce *errors.HttpError
	ur dto.UpdateListItemResponse
	ue *errors.HttpError
	de *errors.HttpError
}

func (f *fakeListSvc) GetListItems(dto.GetListItemsRequest) (dto.GetListItemsResponse, *errors.HttpError) {
	return f.g, f.ge
}

func (f *fakeListSvc) CreateListItem(dto.CreateListItemRequest) (dto.CreateListItemResponse, *errors.HttpError) {
	return f.cr, f.ce
}

func (f *fakeListSvc) UpdateListItem(dto.UpdateListItemRequest) (dto.UpdateListItemResponse, *errors.HttpError) {
	return f.ur, f.ue
}

func (f *fakeListSvc) DeleteListItem(dto.DeleteListItemRequest) *errors.HttpError {
	return f.de
}

func TestGetListItems_Error(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{ge: &errors.HttpError{StatusCode: 404, Message: "m"}})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid")
	ctx.SetParamValues("u1")
	if err := c.getListItems(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 404 {
		t.Fatal(rec.Code)
	}
}

func TestGetListItems_OK(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{g: dto.GetListItemsResponse{}})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid")
	ctx.SetParamValues("u1")
	if err := c.getListItems(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
}

func TestPostListItem_BindError(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid")
	ctx.SetParamValues("u1")
	if err := c.postListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 422 {
		t.Fatal(rec.Code)
	}
}

func TestPostListItem_ServiceError(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{ce: &errors.HttpError{StatusCode: 503, Message: "m"}})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"title":"t"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid")
	ctx.SetParamValues("u1")
	if err := c.postListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 503 {
		t.Fatal(rec.Code)
	}
}

func TestPostListItem_OK(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{cr: dto.CreateListItemResponse{Title: "t"}})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"title":"t"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid")
	ctx.SetParamValues("u1")
	if err := c.postListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
}

func TestPutListItem_BindError(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{})
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader([]byte(`{`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid", "list-item-uuid")
	ctx.SetParamValues("u1", "i1")
	if err := c.putListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 422 {
		t.Fatal(rec.Code)
	}
}

func TestPutListItem_ServiceError(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{ue: &errors.HttpError{StatusCode: 503, Message: "m"}})
	body := `{"title":"t","description":"d","active":true}`
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader([]byte(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid", "list-item-uuid")
	ctx.SetParamValues("u1", "i1")
	if err := c.putListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 503 {
		t.Fatal(rec.Code)
	}
}

func TestPutListItem_OK(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{ur: dto.UpdateListItemResponse{Title: "t"}})
	body := `{"title":"t","description":"d","active":true}`
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader([]byte(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid", "list-item-uuid")
	ctx.SetParamValues("u1", "i1")
	if err := c.putListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
}

func TestDeleteListItem_Error(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{de: &errors.HttpError{StatusCode: 404, Message: "m"}})
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid", "list-item-uuid")
	ctx.SetParamValues("u1", "i1")
	if err := c.deleteListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 404 {
		t.Fatal(rec.Code)
	}
}

func TestDeleteListItem_OK(t *testing.T) {
	e := echo.New()
	c := newListItemController(&fakeListSvc{})
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("user-uuid", "list-item-uuid")
	ctx.SetParamValues("u1", "i1")
	if err := c.deleteListItem(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 204 {
		t.Fatal(rec.Code)
	}
}

func TestTodoInitRoutes(t *testing.T) {
	e := echo.New()
	Init(e, &fakeListSvc{})
	if len(e.Routes()) < 4 {
		t.Fatal(len(e.Routes()))
	}
}
