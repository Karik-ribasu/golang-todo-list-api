package domainServices

import (
	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

type userDomainService struct {
	db data.DbManager
}

func newUserDomainService(db data.DbManager) (dsvc UserDomainService) {
	return &userDomainService{db: db}
}

func (s *userDomainService) GetUserByUUID(userUUID string) (user entity.User, err error) {
	user, err = s.db.UserRepo().GetUserByUUID(userUUID)
	return user, err
}

func (s *userDomainService) GetUserByNickName(nickName string) (user entity.User, err error) {
	user, err = s.db.UserRepo().GetUserByNickName(nickName)
	return user, err
}

func (s *userDomainService) CreateUser(userUUID, nickName string, password []byte) (err error) {
	err = s.db.UserRepo().CreateUser(userUUID, nickName, password)
	return err
}
