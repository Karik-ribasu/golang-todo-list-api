package data

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
)

func TestInitDB_OpenError(t *testing.T) {
	old := openSQL
	defer func() { openSQL = old }()
	openSQL = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("open failed")
	}
	_, err := InitDB(config.Config{Db: config.Db{User: "u", Passwd: "p", Addr: "h", Port: "1", Name: "n"}})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInitDB_PingError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectPing().WillReturnError(errors.New("down"))

	old := openSQL
	defer func() { openSQL = old }()
	openSQL = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	_, err = InitDB(config.Config{Db: config.Db{User: "u", Passwd: "p", Addr: "h", Port: "1", Name: "n"}})
	if err == nil {
		t.Fatal("expected error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestInitDB_OK(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectPing()

	old := openSQL
	defer func() { openSQL = old }()
	openSQL = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}

	dm, err := InitDB(config.Config{Db: config.Db{User: "u", Passwd: "p", Addr: "h", Port: "1", Name: "n"}})
	if err != nil {
		t.Fatal(err)
	}
	if dm == nil {
		t.Fatal("nil manager")
	}
	if dm.ListItemRepo() == nil || dm.UserRepo() == nil {
		t.Fatal("repos")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNewManagerFromDB(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	dm := NewManagerFromDB(db)
	if dm.ListItemRepo() == nil {
		t.Fatal("expected list repo")
	}
}
