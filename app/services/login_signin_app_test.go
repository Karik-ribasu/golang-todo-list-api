package appServices

import (
	"crypto/rsa"
	"database/sql"
	"errors"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser_Success(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	u := &stubUserDom{u: entity.User{UserUUID: "uuid-1", Password: hash}}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	resp, herr := svc.LoginUser(dto.LoginRequest{NickName: "n", Password: "secret"})
	if herr != nil || resp.Token == "" {
		t.Fatalf("%v %+v", herr, resp)
	}
}

func TestLoginUser_UserNotFound(t *testing.T) {
	u := &stubUserDom{err: sql.ErrNoRows}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	_, herr := svc.LoginUser(dto.LoginRequest{NickName: "n", Password: "p"})
	if herr == nil || herr.StatusCode != 404 {
		t.Fatal()
	}
}

func TestLoginUser_OtherSQL(t *testing.T) {
	u := &stubUserDom{err: errors.New("boom")}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	_, herr := svc.LoginUser(dto.LoginRequest{NickName: "n", Password: "p"})
	if herr == nil || herr.StatusCode != 503 {
		t.Fatal()
	}
}

func TestLoginUser_BadPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("ok"), bcrypt.MinCost)
	u := &stubUserDom{u: entity.User{Password: hash}}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	_, herr := svc.LoginUser(dto.LoginRequest{NickName: "n", Password: "wrong"})
	if herr == nil || herr.StatusCode != 401 {
		t.Fatal()
	}
}

func TestSiginUser_CreateError(t *testing.T) {
	u := &stubUserDom{ce: sql.ErrNoRows}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	if herr := svc.SiginUser(dto.SigninRequest{NickName: "n", Password: "p"}); herr == nil {
		t.Fatal()
	}
}

func TestSiginUser_GenericSQLError(t *testing.T) {
	u := &stubUserDom{ce: errors.New("other")}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	if herr := svc.SiginUser(dto.SigninRequest{NickName: "n", Password: "p"}); herr == nil || herr.StatusCode != 503 {
		t.Fatal(herr)
	}
}

func TestSiginUser_OK(t *testing.T) {
	u := &stubUserDom{}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	if herr := svc.SiginUser(dto.SigninRequest{NickName: "n", Password: "p"}); herr != nil {
		t.Fatal(herr)
	}
}

func TestLoginUser_JWTError(t *testing.T) {
	prev := jwtTokenForUser
	defer func() { jwtTokenForUser = prev }()
	jwtTokenForUser = func(*rsa.PrivateKey, string) (auth.AuthToken, error) {
		return auth.AuthToken{}, errors.New("jwt")
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	u := &stubUserDom{u: entity.User{UserUUID: "id", Password: hash}}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	_, herr := svc.LoginUser(dto.LoginRequest{NickName: "n", Password: "secret"})
	if herr == nil || herr.StatusCode != 500 {
		t.Fatal(herr)
	}
}

func TestSiginUser_BcryptError(t *testing.T) {
	prev := bcryptFromPassword
	defer func() { bcryptFromPassword = prev }()
	bcryptFromPassword = func([]byte, int) ([]byte, error) { return nil, errors.New("bcrypt") }
	u := &stubUserDom{}
	svc := newLoginSiginAppService(rsaConfig(t), u)
	if herr := svc.SiginUser(dto.SigninRequest{NickName: "n", Password: "p"}); herr == nil || herr.StatusCode != 500 {
		t.Fatal(herr)
	}
}
