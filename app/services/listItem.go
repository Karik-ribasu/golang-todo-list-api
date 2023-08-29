package appServices

import (
	"encoding/json"
	"net/http"

	"github.com/Karik-ribasu/golang-todo-list-api/app/dto"
	domainServices "github.com/Karik-ribasu/golang-todo-list-api/domain/services"
	"github.com/Karik-ribasu/golang-todo-list-api/infra/errors"
	"github.com/google/uuid"
)

type listItemAppService struct {
	listItemDomainService domainServices.ListItemDomainService
	userDomainService     domainServices.UserDomainService
}

func newListItemAppService(listItemDomainService domainServices.ListItemDomainService, userDomainService domainServices.UserDomainService) ListItemAppService {
	return &listItemAppService{listItemDomainService: listItemDomainService}
}

func (s *listItemAppService) GetListItems(reqData dto.GetListItemsRequest) (resp dto.GetListItemsResponse, err *errors.HttpError) {
	if reqData.UserUUID == "" {
		return resp, &errors.HttpError{StatusCode: http.StatusNotFound, Message: `{"error": "Not Found"}`}
	}

	user, errSQL := s.userDomainService.GetUserByUUID(reqData.UserUUID)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return resp, httpError
	}

	listItem, errSQL := s.listItemDomainService.GetListItems(user.UserID)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return resp, httpError
	}

	listItemData, errJSON := json.Marshal(listItem)
	if errJSON != nil {
		return resp, &errors.HttpError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	errJSON = json.Unmarshal(listItemData, &resp)
	if errJSON != nil {
		return resp, &errors.HttpError{StatusCode: http.StatusInternalServerError, Message: "Internal Server Error"}
	}

	return
}

func (s *listItemAppService) CreateListItem(reqData dto.CreateListItemRequest) (resp dto.CreateListItemResponse, err *errors.HttpError) {
	if reqData.UserUUID == "" {
		return resp, &errors.HttpError{StatusCode: http.StatusNotFound, Message: `{"error": "Not Found"}`}
	}

	user, errSQL := s.userDomainService.GetUserByUUID(reqData.UserUUID)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return resp, httpError
	}

	listItemUUID := uuid.New().String()
	errSQL = s.listItemDomainService.CreateListItem(user.UserID, listItemUUID, reqData.Title, reqData.Description)
	if err != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return resp, httpError
	}

	resp.ListItemUUID = listItemUUID
	resp.Title = reqData.Title
	resp.Description = reqData.Description
	resp.Active = true

	return
}

func (s *listItemAppService) UpdateListItem(reqData dto.UpdateListItemRequest) (resp dto.UpdateListItemResponse, err *errors.HttpError) {
	if reqData.UserUUID == "" || reqData.ListItemUUID == "" {
		return resp, &errors.HttpError{StatusCode: http.StatusNotFound, Message: `{"error": "Not Found"}`}
	}

	user, errSQL := s.userDomainService.GetUserByUUID(reqData.UserUUID)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return resp, httpError
	}

	errSQL = s.listItemDomainService.UpdateListItem(user.UserID, reqData.ListItemUUID, reqData.Title, reqData.Description, reqData.Active)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return resp, httpError
	}

	resp = dto.UpdateListItemResponse{
		ListItemUUID: reqData.ListItemUUID,
		Title:        reqData.Title,
		Description:  reqData.Description,
		Active:       reqData.Active,
	}

	return
}

func (s *listItemAppService) DeleteListItem(reqData dto.DeleteListItemRequest) (err *errors.HttpError) {
	if reqData.UserUUID == "" || reqData.ListItemUUID == "" {
		return &errors.HttpError{StatusCode: http.StatusNotFound, Message: `{"error": "Not Found"}`}
	}

	user, errSQL := s.userDomainService.GetUserByUUID(reqData.UserUUID)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return httpError
	}

	errSQL = s.listItemDomainService.DeleteListItem(user.UserID, reqData.ListItemUUID)
	if errSQL != nil {
		httpError := errors.SQLErrorCheck(errSQL)
		return httpError
	}

	return
}
