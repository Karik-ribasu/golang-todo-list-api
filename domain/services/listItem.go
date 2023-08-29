package domainServices

import (
	"github.com/Karik-ribasu/golang-todo-list-api/domain/data"
	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

type listItemDomainService struct {
	listItemRepo data.ListItemRepo
}

func newListItemDomainService(db data.DbManager) (dsvc ListItemDomainService) {
	return &listItemDomainService{listItemRepo: db.ListItemRepo()}
}

func (s *listItemDomainService) GetListItems(userID int64) (listItems []entity.ListItem, err error) {
	listItems, err = s.listItemRepo.GetListItems(userID)
	return listItems, err
}

func (s *listItemDomainService) CreateListItem(userID int64, listItemUUID string, title string, description string) (err error) {
	err = s.listItemRepo.CreateListItem(userID, listItemUUID, title, description)
	return
}

func (s *listItemDomainService) UpdateListItem(userID int64, listItemUUID string, title string, description string, active bool) (err error) {
	err = s.listItemRepo.UpdateListItem(userID, listItemUUID, title, description, active)
	return
}

func (s *listItemDomainService) DeleteListItem(userID int64, listItemUUID string) (err error) {
	err = s.listItemRepo.DeleteListItem(userID, listItemUUID)
	return
}
