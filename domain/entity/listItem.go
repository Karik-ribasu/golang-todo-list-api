package entity

type ListItem struct {
	ListItemID   int64  `json:"-"`
	UserID       int64  `json:"-"`
	ListItemUUID string `json:"list_item_uuid"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
}
