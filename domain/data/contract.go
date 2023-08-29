package data

import "github.com/Karik-ribasu/golang-todo-list-api/domain/entity"

type ListItemRepo interface {
	CreateListItem(userID int64, listItemUUID, title, description string) (err error)
	GetListItems(userID int64) (listItems []entity.ListItem, err error)
	UpdateListItem(userID int64, listItemUUID, title, description string, active bool) (err error)
	DeleteListItem(userID int64, listItemUUID string) (err error)
}

type UserRepo interface {
	CreateUser(userUUID, nickName string, password []byte) (err error)
	GetUserByUUID(userUUID string) (user entity.User, err error)
	GetUserByNickName(nickName string) (user entity.User, err error)
}

type DbManager interface {
	ListItemRepo() ListItemRepo
	UserRepo() UserRepo
}

func (conn *conn) ListItemRepo() ListItemRepo {
	return &listItemRepo{db: conn.db}
}

func (conn *conn) UserRepo() UserRepo {
	return &userRepo{db: conn.db}
}
