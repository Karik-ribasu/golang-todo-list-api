package domainServices

import (
	"errors"
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

type stubUserRepo struct {
	u    entity.User
	err  error
	cerr error
}

func (s *stubUserRepo) CreateUser(userUUID, nickName string, password []byte) error {
	return s.cerr
}

func (s *stubUserRepo) GetUserByUUID(userUUID string) (entity.User, error) {
	return s.u, s.err
}

func (s *stubUserRepo) GetUserByNickName(nickName string) (entity.User, error) {
	return s.u, s.err
}

type stubListRepo struct {
	items []entity.ListItem
	err   error
	cerr  error
	uerr  error
	derr  error
}

func (s *stubListRepo) CreateListItem(userID int64, listItemUUID, title, description string) error {
	return s.cerr
}

func (s *stubListRepo) GetListItems(userID int64) ([]entity.ListItem, error) {
	return s.items, s.err
}

func (s *stubListRepo) UpdateListItem(userID int64, listItemUUID, title, description string, active bool) error {
	return s.uerr
}

func (s *stubListRepo) DeleteListItem(userID int64, listItemUUID string) error {
	return s.derr
}

type stubDB struct {
	u *stubUserRepo
	l *stubListRepo
}

func (s *stubDB) ListItemRepo() data.ListItemRepo { return s.l }
func (s *stubDB) UserRepo() data.UserRepo         { return s.u }

func TestNewDomainSVCAndAccessors(t *testing.T) {
	u := &stubUserRepo{}
	l := &stubListRepo{}
	db := &stubDB{u: u, l: l}
	m := NewDomainSVC(db)
	if m.ListItemDomainService() == nil || m.UserDomainService() == nil {
		t.Fatal()
	}
}

func TestUserDomainService(t *testing.T) {
	u := &stubUserRepo{u: entity.User{UserID: 5}}
	svc := newUserDomainService(&stubDB{u: u})
	got, err := svc.GetUserByUUID("x")
	if err != nil || got.UserID != 5 {
		t.Fatal()
	}
	u.err = errors.New("e")
	if _, err := svc.GetUserByNickName("n"); err == nil {
		t.Fatal()
	}
	u.err = nil
	u.cerr = errors.New("c")
	if err := svc.CreateUser("a", "b", nil); err == nil {
		t.Fatal()
	}
}

func TestListItemDomainService(t *testing.T) {
	l := &stubListRepo{items: []entity.ListItem{{ListItemUUID: "z"}}}
	svc := newListItemDomainService(&stubDB{l: l})
	got, err := svc.GetListItems(1)
	if err != nil || len(got) != 1 {
		t.Fatal()
	}
	l.err = errors.New("e")
	if _, err := svc.GetListItems(1); err == nil {
		t.Fatal()
	}
	l.err = nil
	l.cerr = errors.New("c")
	if err := svc.CreateListItem(1, "u", "t", "d"); err == nil {
		t.Fatal()
	}
	l.cerr = nil
	l.uerr = errors.New("u")
	if err := svc.UpdateListItem(1, "u", "t", "d", true); err == nil {
		t.Fatal()
	}
	l.uerr = nil
	l.derr = errors.New("d")
	if err := svc.DeleteListItem(1, "u"); err == nil {
		t.Fatal()
	}
}
