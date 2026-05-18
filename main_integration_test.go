//go:build integration

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func TestIntegration_TodoFlow(t *testing.T) {
	root := findRepoRoot(t)
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		t.Fatal(err)
	}
	pemBytes, err := config.LoadCertificatePEM(cfg.App)
	if err != nil {
		t.Fatal(err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		t.Fatal(err)
	}
	cfg.App.PrivateKey = key

	dm, err := data.InitDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	e := InitializeServer(cfg, dm)

	nick := fmt.Sprintf("it%d", time.Now().UnixNano())
	password := "p@ssw0rd!"

	rec := doJSON(t, e, http.MethodPost, "/sign-in", "", map[string]string{"nick_name": nick, "password": password})
	if rec.Code != http.StatusNoContent {
		t.Fatalf("sign-in: %d %s", rec.Code, rec.Body.String())
	}

	token := loginToken(t, e, nick, password)
	userUUID := userUUIDFromJWT(t, token)

	rec = doJSON(t, e, http.MethodPost, "/user/"+userUUID+"/todo-list", token, map[string]string{"title": "buy milk", "description": "full fat"})
	if rec.Code != http.StatusOK {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	itemUUID, _ := created["list_item_uuid"].(string)
	if itemUUID == "" {
		t.Fatalf("missing list_item_uuid: %#v", created)
	}

	rec = doJSON(t, e, http.MethodGet, "/user/"+userUUID+"/todo-list", token, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list: %d %s", rec.Code, rec.Body.String())
	}
	var items []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %#v", items)
	}

	rec = doJSON(t, e, http.MethodPut, "/user/"+userUUID+"/todo-list/"+itemUUID, token, map[string]any{
		"title": "buy oat milk", "description": "barista", "active": false,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("update: %d %s", rec.Code, rec.Body.String())
	}

	rec = doJSON(t, e, http.MethodDelete, "/user/"+userUUID+"/todo-list/"+itemUUID, token, nil)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: %d %s", rec.Code, rec.Body.String())
	}

	rec = doJSON(t, e, http.MethodGet, "/user/"+userUUID+"/todo-list", token, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list after delete: %d %s", rec.Code, rec.Body.String())
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty list, got %#v", items)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("go.mod not found from %s", wd)
	return ""
}

func doJSON(t *testing.T, e http.Handler, method, path, bearer string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reqBody = strings.NewReader(string(b))
	}
	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	if bearer != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func loginToken(t *testing.T, e http.Handler, nick, password string) string {
	t.Helper()
	rec := doJSON(t, e, http.MethodPost, "/log-in", "", map[string]string{"nick_name": nick, "password": password})
	if rec.Code != http.StatusOK {
		t.Fatalf("log-in: %d %s", rec.Code, rec.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	tok, _ := out["token"].(string)
	if tok == "" {
		t.Fatalf("no token in %#v", out)
	}
	return tok
}

func userUUIDFromJWT(t *testing.T, token string) string {
	t.Helper()
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		t.Fatal("invalid jwt")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(payload, &m); err != nil {
		t.Fatal(err)
	}
	u, _ := m["user_uuid"].(string)
	if u == "" {
		t.Fatalf("user_uuid missing in %s", string(payload))
	}
	return u
}
