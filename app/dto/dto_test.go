package dto

import (
	"encoding/json"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/infra/auth"
)

func TestLoginResponseJSON(t *testing.T) {
	r := LoginResponse{
		AuthToken: auth.AuthToken{Token: "t", IssuedAt: 1, ExpiresAt: 2},
		Message:   "ok",
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m["token"] != "t" {
		t.Fatalf("%s", b)
	}
}

func TestSigninRequestRoundTrip(t *testing.T) {
	s := SigninRequest{NickName: "n", Password: "p"}
	b, _ := json.Marshal(s)
	var out SigninRequest
	json.Unmarshal(b, &out)
	if out.NickName != "n" {
		t.Fatal()
	}
}
