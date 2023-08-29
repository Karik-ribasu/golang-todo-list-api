package entity

type User struct {
	UserID   int64
	UserUUID string
	NickName string
	Password []byte
}
