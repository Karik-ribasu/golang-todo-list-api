package signInRoute

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/errors"
	"github.com/labstack/echo/v4"
)

type fakeLoginSvc struct {
	sig  *errors.HttpError
	resp dto.LoginResponse
	le   *errors.HttpError
}

func (f *fakeLoginSvc) LoginUser(dto.LoginRequest) (dto.LoginResponse, *errors.HttpError) {
	return f.resp, f.le
}

func (f *fakeLoginSvc) SiginUser(dto.SigninRequest) *errors.HttpError {
	return f.sig
}

func TestHandleUserSignIn_BindError(t *testing.T) {
	e := echo.New()
	c := newLoginSiginController(&fakeLoginSvc{})
	req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader([]byte(`{`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	if err := c.handleUserSignIn(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 422 {
		t.Fatal(rec.Code)
	}
}

func TestHandleUserSignIn_ServiceError(t *testing.T) {
	e := echo.New()
	c := newLoginSiginController(&fakeLoginSvc{sig: &errors.HttpError{StatusCode: 503, Message: "x"}})
	req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader([]byte(`{"nick_name":"n","password":"p"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	if err := c.handleUserSignIn(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 503 {
		t.Fatal(rec.Code)
	}
}

func TestHandleUserSignIn_OK(t *testing.T) {
	e := echo.New()
	c := newLoginSiginController(&fakeLoginSvc{})
	req := httptest.NewRequest(http.MethodPost, "/sign-in", bytes.NewReader([]byte(`{"nick_name":"n","password":"p"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	if err := c.handleUserSignIn(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 204 {
		t.Fatal(rec.Code)
	}
}

func TestHandleUserLogin_BindError(t *testing.T) {
	e := echo.New()
	c := newLoginSiginController(&fakeLoginSvc{})
	req := httptest.NewRequest(http.MethodPost, "/log-in", bytes.NewReader([]byte(`{`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	if err := c.handleUserLogin(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 422 {
		t.Fatal(rec.Code)
	}
}

func TestHandleUserLogin_ServiceError(t *testing.T) {
	e := echo.New()
	c := newLoginSiginController(&fakeLoginSvc{le: &errors.HttpError{StatusCode: 401, Message: "n"}})
	req := httptest.NewRequest(http.MethodPost, "/log-in", bytes.NewReader([]byte(`{"nick_name":"n","password":"p"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	if err := c.handleUserLogin(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 401 {
		t.Fatal(rec.Code)
	}
}

func TestHandleUserLogin_OK(t *testing.T) {
	e := echo.New()
	c := newLoginSiginController(&fakeLoginSvc{resp: dto.LoginResponse{Message: "ok"}})
	req := httptest.NewRequest(http.MethodPost, "/log-in", bytes.NewReader([]byte(`{"nick_name":"n","password":"p"}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	if err := c.handleUserLogin(ctx); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 {
		t.Fatal(rec.Code)
	}
}

func TestInitRegistersRoutes(t *testing.T) {
	e := echo.New()
	Init(e, &fakeLoginSvc{})
	routes := e.Routes()
	if len(routes) < 2 {
		t.Fatal(len(routes))
	}
}
