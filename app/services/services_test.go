package appServices

import (
	"testing"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
	domainServices "github.com/Karik-ribasu/golang-todo-list-api/domain/services"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
)

type stubUserDom struct {
	u   entity.User
	err error
	ce  error
}

func (s *stubUserDom) GetUserByUUID(string) (entity.User, error)     { return s.u, s.err }
func (s *stubUserDom) GetUserByNickName(string) (entity.User, error) { return s.u, s.err }
func (s *stubUserDom) CreateUser(string, string, []byte) error       { return s.ce }

type stubListDom struct {
	items []entity.ListItem
	ge    error
	ce    error
	ue    error
	de    error
}

func (s *stubListDom) GetListItems(int64) ([]entity.ListItem, error)            { return s.items, s.ge }
func (s *stubListDom) CreateListItem(int64, string, string, string) error       { return s.ce }
func (s *stubListDom) UpdateListItem(int64, string, string, string, bool) error { return s.ue }
func (s *stubListDom) DeleteListItem(int64, string) error                       { return s.de }

type stubDomMgr struct {
	u *stubUserDom
	l *stubListDom
}

func (s *stubDomMgr) UserDomainService() domainServices.UserDomainService         { return s.u }
func (s *stubDomMgr) ListItemDomainService() domainServices.ListItemDomainService { return s.l }

func TestNewAppServiceManagers(t *testing.T) {
	u := &stubUserDom{}
	l := &stubListDom{}
	m := NewAppService(config.Config{}, &stubDomMgr{u: u, l: l})
	if m.TodoListAppService() == nil || m.LoginSiginAppService() == nil {
		t.Fatal()
	}
}
