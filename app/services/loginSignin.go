package appServices

import (
	"net/http"

	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	domainServices "github.com/Karik-ribasu/golang-todo-list-api/domain/services"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/auth"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type loginSiginAppService struct {
	userDomainService domainServices.UserDomainService
	cfg               config.Config
}

func newLoginSiginAppService(cfg config.Config, userDomainService domainServices.UserDomainService) LoginSiginAppService {
	return &loginSiginAppService{
		userDomainService: userDomainService,
		cfg:               cfg,
	}
}

func (s *loginSiginAppService) LoginUser(reqData dto.LoginRequest) (resp dto.LoginResponse, errHttp *errors.HttpError) {

	user, err := s.userDomainService.GetUserByNickName(reqData.NickName)
	if err != nil {
		errHttp = errors.SQLErrorCheck(err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password))
	if err != nil {
		errHttp.StatusCode = http.StatusUnauthorized
		errHttp.Message = "NickName or Password invalid"
		return
	}

	authToken, err := auth.GenerateJWTToken(s.cfg.App.PrivateKey, user.UserUUID)
	if err != nil {
		return resp, &errors.HttpError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	resp.AuthToken = authToken

	return
}

func (s *loginSiginAppService) SiginUser(reqData dto.SigninRequest) (errHttp *errors.HttpError) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), 1)
	if err != nil {
		return &errors.HttpError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	userUUID := uuid.New().String()
	err = s.userDomainService.CreateUser(userUUID, reqData.NickName, hashedPassword)
	if err != nil {
		errHttp = errors.SQLErrorCheck(err)
		return
	}

	return
}
