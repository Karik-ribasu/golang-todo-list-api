package data

import (
	"database/sql"

	"github.com/Karik-ribasu/golang-todo-list-api/domain/entity"
)

type listItemRepo struct {
	db *sql.DB
}

func (l listItemRepo) GetListItems(userID int64) (listItems []entity.ListItem, err error) {
	query := `
	Select 
		li.list_item_id,
		li.list_item_uuid,
		li.user_id,
		li.title,
		li.description,
		li.active
	From list_item li
	Where li.user_id = ?
	`
	result, err := l.db.Query(query, userID)
	if err != nil {
		return listItems, err
	}

	for result.Next() {
		listItem := entity.ListItem{}
		err = result.Scan(&listItem.ListItemID, &listItem.ListItemUUID, &listItem.UserID, &listItem.Title, &listItem.Description, &listItem.Active)
		if err != nil {
			return listItems, err
		}

		listItems = append(listItems, listItem)
	}

	return listItems, err
}

func (l listItemRepo) CreateListItem(userID int64, listItemUUID, title, description string) (err error) {
	query := `
	Insert into list_item(list_item_uuid, user_id, title, description)
	Values(?, ?, ?, ?)
	`

	_, err = l.db.Exec(query, listItemUUID, userID, title, description)
	return err
}

func (l listItemRepo) UpdateListItem(userID int64, listItemUUID, title, description string, active bool) (err error) {
	query := `
	update list_item
	set
		title = ?,
		description = ?,
		active = ?
	where
		user_id = ? and list_item_uuid = ?
	`

	_, err = l.db.Exec(query, title, description, active, userID, listItemUUID)
	return err
}

func (l listItemRepo) DeleteListItem(userID int64, listItemUUID string) (err error) {
	query := `
	delete from list_item
	where
		user_id = ? and list_item_uuid = ?
	`

	_, err = l.db.Exec(query, userID, listItemUUID)
	return err
}
