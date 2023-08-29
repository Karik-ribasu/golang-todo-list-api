package appServices

import (
	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	domainServices "github.com/Karik-ribasu/golang-todo-list-api/domain/services"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/config"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/errors"
)

type appServiceManager struct {
	domainSvcManager domainServices.DomainSvcManager
	cfg              config.Config
}

type AppServiceManager interface {
	TodoListAppService() ListItemAppService
	LoginSiginAppService() LoginSiginAppService
}

type ListItemAppService interface {
	GetListItems(dto.GetListItemsRequest) (dto.GetListItemsResponse, *errors.HttpError)
	CreateListItem(dto.CreateListItemRequest) (dto.CreateListItemResponse, *errors.HttpError)
	UpdateListItem(dto.UpdateListItemRequest) (dto.UpdateListItemResponse, *errors.HttpError)
	DeleteListItem(dto.DeleteListItemRequest) *errors.HttpError
}

type LoginSiginAppService interface {
	LoginUser(dto.LoginRequest) (resp dto.LoginResponse, errHttp *errors.HttpError)
	SiginUser(reqData dto.SigninRequest) *errors.HttpError
}

func NewAppService(cfg config.Config, domainSvcManager domainServices.DomainSvcManager) AppServiceManager {
	return &appServiceManager{
		cfg:              cfg,
		domainSvcManager: domainSvcManager,
	}
}

func (svc *appServiceManager) TodoListAppService() ListItemAppService {
	return newListItemAppService(svc.domainSvcManager.ListItemDomainService(), svc.domainSvcManager.UserDomainService())
}

func (svc *appServiceManager) LoginSiginAppService() LoginSiginAppService {
	return newLoginSiginAppService(svc.cfg, svc.domainSvcManager.UserDomainService())
}
