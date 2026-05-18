package data

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepo_CreateUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := userRepo{db: db}
	mock.ExpectExec("Insert into user").WillReturnError(errors.New("e"))
	if err := repo.CreateUser("u", "n", []byte("p")); err == nil {
		t.Fatal("expected error")
	}
}

func TestUserRepo_CreateUser_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := userRepo{db: db}
	mock.ExpectExec("Insert into user").WithArgs("uuid", "nick", []byte("pw")).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.CreateUser("uuid", "nick", []byte("pw")); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_GetUserByUUID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := userRepo{db: db}
	mock.ExpectQuery("Select").WillReturnError(errors.New("e"))
	_, err = repo.GetUserByUUID("uuid")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUserRepo_GetUserByUUID_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := userRepo{db: db}
	rows := sqlmock.NewRows([]string{"user_id", "user_uuid", "nick_name", "password"}).AddRow(int64(1), "uuid", "nick", []byte("pw"))
	mock.ExpectQuery("Select").WithArgs("uuid").WillReturnRows(rows)
	u, err := repo.GetUserByUUID("uuid")
	if err != nil {
		t.Fatal(err)
	}
	if u.UserID != 1 || u.NickName != "nick" {
		t.Fatalf("%+v", u)
	}
}

func TestUserRepo_GetUserByNickName_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := userRepo{db: db}
	mock.ExpectQuery("Select").WillReturnError(errors.New("e"))
	_, err = repo.GetUserByNickName("nick")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUserRepo_GetUserByNickName_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := userRepo{db: db}
	rows := sqlmock.NewRows([]string{"user_id", "user_uuid", "nick_name", "password"}).AddRow(int64(1), "uuid", "nick", []byte("pw"))
	mock.ExpectQuery("Select").WithArgs("nick").WillReturnRows(rows)
	u, err := repo.GetUserByNickName("nick")
	if err != nil {
		t.Fatal(err)
	}
	if u.UserUUID != "uuid" {
		t.Fatalf("%+v", u)
	}
}
