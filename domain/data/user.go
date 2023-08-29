package data

import (
	"database/sql"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

type userRepo struct {
	db *sql.DB
}

func (l userRepo) CreateUser(userUUID, nickName string, password []byte) (err error) {
	query := `
	Insert into user(user_uuid, nick_name, password)
	Values(?, ?, ?)
	`
	_, err = l.db.Exec(query, userUUID, nickName, password)

	return err
}

func (l userRepo) GetUserByUUID(userUUID string) (user entity.User, err error) {
	query := `
	Select
		user_id,
		user_uuid,
		nick_name,
		password
	From user
	Where user_uuid = ? and active = true
	`
	result := l.db.QueryRow(query, userUUID)

	err = result.Scan(&user.UserID, &user.UserUUID, &user.NickName, &user.Password)
	return user, err
}

func (l userRepo) GetUserByNickName(nickName string) (user entity.User, err error) {
	query := `
	Select
		user_id,
		user_uuid,
		nick_name,
		password
	From user
	Where nick_name = ? and active = true
	`
	result := l.db.QueryRow(query, nickName)

	err = result.Scan(&user.UserID, &user.UserUUID, &user.NickName, &user.Password)
	return user, err
}
