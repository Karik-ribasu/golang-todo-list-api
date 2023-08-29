package dto

import "github.com/Karik-ribasu/golang-todo-list-api/infra/auth"

type LoginRequest struct {
	NickName string `json:"nick_name"`
	Password string `json:"password"`
}

type SigninRequest struct {
	NickName string `json:"nick_name"`
	Password string `json:"password"`
}

type LoginResponse struct {
	auth.AuthToken
	Message string `json:"message"`
}

type SigninResponse struct {
	Message string `json:"message"`
}
