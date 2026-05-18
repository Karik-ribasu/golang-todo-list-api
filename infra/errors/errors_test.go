package errors

import (
	"database/sql"
	"errors"
	"testing"
)

func TestSQLErrorCheck_NoRows(t *testing.T) {
	h := SQLErrorCheck(sql.ErrNoRows)
	if h.StatusCode != 404 {
		t.Fatal()
	}
}

func TestSQLErrorCheck_Other(t *testing.T) {
	h := SQLErrorCheck(errors.New("x"))
	if h.StatusCode != 503 {
		t.Fatal()
	}
}
