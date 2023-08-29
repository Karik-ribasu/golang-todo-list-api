package domainServices

import (
	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

type domainSvcManager struct {
	db data.DbManager
}

type DomainSvcManager interface {
	ListItemDomainService() ListItemDomainService
	UserDomainService() UserDomainService
}

type UserDomainService interface {
	GetUserByUUID(userUUID string) (user entity.User, err error)
	GetUserByNickName(nickName string) (user entity.User, err error)
	CreateUser(userUUID, nickName string, password []byte) (err error)
}

type ListItemDomainService interface {
	GetListItems(userID int64) (listItems []entity.ListItem, err error)
	CreateListItem(userID int64, listItemUUID string, title string, description string) (err error)
	UpdateListItem(userID int64, listItemUUID string, title string, description string, active bool) (err error)
	DeleteListItem(userID int64, listItemUUID string) (err error)
}

func NewDomainSVC(db data.DbManager) DomainSvcManager {
	return &domainSvcManager{db: db}
}

func (svc *domainSvcManager) ListItemDomainService() ListItemDomainService {
	return newListItemDomainService(svc.db)
}

func (svc *domainSvcManager) UserDomainService() UserDomainService {
	return newUserDomainService(svc.db)
}
