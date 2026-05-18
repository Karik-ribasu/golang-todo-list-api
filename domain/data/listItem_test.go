package data

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestListItemRepo_GetListItems_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectQuery("Select").WillReturnError(errors.New("boom"))
	_, err = repo.GetListItems(1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListItemRepo_GetListItems_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	rows := sqlmock.NewRows([]string{"a"}).AddRow("bad")
	mock.ExpectQuery("Select").WillReturnRows(rows)
	_, err = repo.GetListItems(1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListItemRepo_GetListItems_RowsErr(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	rows := sqlmock.NewRows([]string{"list_item_id", "list_item_uuid", "user_id", "title", "description", "active"}).
		AddRow(1, "uuid", 2, "t", "d", true).
		RowError(0, errors.New("row err"))
	mock.ExpectQuery("Select").WillReturnRows(rows)
	_, err = repo.GetListItems(1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListItemRepo_GetListItems_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	rows := sqlmock.NewRows([]string{"list_item_id", "list_item_uuid", "user_id", "title", "description", "active"}).
		AddRow(1, "u1", 2, "t", "d", true)
	mock.ExpectQuery("Select").WithArgs(int64(2)).WillReturnRows(rows)
	got, err := repo.GetListItems(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ListItemUUID != "u1" {
		t.Fatalf("%#v", got)
	}
}

func TestListItemRepo_CreateListItem_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectExec("Insert into list_item").WillReturnError(errors.New("e"))
	if err := repo.CreateListItem(1, "uuid", "t", "d"); err == nil {
		t.Fatal("expected error")
	}
}

func TestListItemRepo_CreateListItem_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectExec("Insert into list_item").WithArgs("uuid", int64(1), "t", "d").WillReturnResult(sqlmock.NewResult(1, 1))
	if err := repo.CreateListItem(1, "uuid", "t", "d"); err != nil {
		t.Fatal(err)
	}
}

func TestListItemRepo_UpdateListItem_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectExec("update list_item").WillReturnError(errors.New("e"))
	if err := repo.UpdateListItem(1, "uuid", "t", "d", true); err == nil {
		t.Fatal("expected error")
	}
}

func TestListItemRepo_UpdateListItem_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectExec("update list_item").WithArgs("t", "d", true, int64(1), "uuid").WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.UpdateListItem(1, "uuid", "t", "d", true); err != nil {
		t.Fatal(err)
	}
}

func TestListItemRepo_DeleteListItem_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectExec("delete from list_item").WillReturnError(errors.New("e"))
	if err := repo.DeleteListItem(1, "uuid"); err == nil {
		t.Fatal("expected error")
	}
}

func TestListItemRepo_DeleteListItem_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := listItemRepo{db: db}
	mock.ExpectExec("delete from list_item").WithArgs(int64(1), "uuid").WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repo.DeleteListItem(1, "uuid"); err != nil {
		t.Fatal(err)
	}
}

func TestConnRepos(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	c := &conn{db: db}
	var _ ListItemRepo = c.ListItemRepo()
	var _ UserRepo = c.UserRepo()
}
